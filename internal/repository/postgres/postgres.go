package postgres

import (
	"avito/internal/logger"
	"avito/internal/models"
	"avito/internal/router/handlers/responses"
	"context"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepositoryImpl struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewRepositoryImpl(dsn string, logger *logger.Logger) (*PostgresRepositoryImpl, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresRepositoryImpl{
		logger: logger,
		db:     db,
	}, nil
}

func (p *PostgresRepositoryImpl) CreateUser(ctx context.Context,
	userName string,
	hashedPassword []byte,
	tokenString string) error {

	tx := p.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &models.User{
		Username:       userName,
		HashedPassword: hashedPassword,
		Token:          tokenString,
		Balance:        1000,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *PostgresRepositoryImpl) CreateTransaction(ctx context.Context,
	senderName string,
	receiverName string,
	amount int) error {

	tx := p.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var sender, receiver models.User

	if err := tx.Where("username = ?", senderName).First(&sender).Error; err != nil {
		tx.Rollback()
		return err
	}

	if sender.Balance < amount {
		tx.Rollback()
		return errors.New("недостаточно средств")
	}

	if err := tx.Where("username = ?", receiverName).First(&receiver).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&sender).Update("balance", sender.Balance-amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&receiver).Update("balance", receiver.Balance+amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	transaction := &models.Transaction{
		Amount:     amount,
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *PostgresRepositoryImpl) CreatePurchase(ctx context.Context,
	purchaserName string, merchName string) error {

	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var purchaser models.User
	if err := tx.Where("username = ?", purchaserName).First(&purchaser).Error; err != nil {
		tx.Rollback()
		return err
	}

	var merch models.Merch
	if err := tx.Where("name = ?", merchName).First(&merch).Error; err != nil {
		tx.Rollback()
		return err
	}

	if purchaser.Balance < merch.Price {
		tx.Rollback()
		return errors.New("insufficient funds for purchase")
	}

	if err := tx.Model(&purchaser).Update("balance", purchaser.Balance-merch.Price).Error; err != nil {
		tx.Rollback()
		return err
	}

	purchase := &models.Purchase{
		UserID:  purchaser.ID,
		MerchID: merch.ID,
		Price:   merch.Price,
	}

	if err := tx.Create(purchase).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *PostgresRepositoryImpl) FindUserByName(ctx context.Context, userName string) (*models.User, error) {
	var user models.User
	if err := p.db.WithContext(ctx).Where("username = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepositoryImpl) FindAppliedTransactions(ctx context.Context, sentORreceived bool, userName string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	user, err := p.FindUserByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	if sentORreceived {
		if err := p.db.WithContext(ctx).
			Where("senderid = ?", user.ID).
			Find(&transactions).Error; err != nil {
			return nil, err
		}
	} else {
		if err := p.db.WithContext(ctx).
			Where("receiverid = ?", user.ID).
			Find(&transactions).Error; err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func (p *PostgresRepositoryImpl) FindBoughtMerch(ctx context.Context, userName string) ([]responses.InventoryItem, error) {
	var results []responses.InventoryItem

	user, err := p.FindUserByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	if err = p.db.WithContext(ctx).
		Table("purchases").
		Select("merch.name AS merch_name, COUNT(*) AS count").
		Joins("JOIN merch ON purchases.merch_id = merch.id").
		Where("purchases.user_id = ?", user.ID).
		Group("merch.name").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
