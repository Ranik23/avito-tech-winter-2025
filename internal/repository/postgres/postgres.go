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

func (p *PostgresRepositoryImpl) FindTransactions(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := p.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (p *PostgresRepositoryImpl) CreateUser(ctx context.Context, userName string, hashedPassword []byte, tokenString string) error {
	user := &models.User{
		Username:       userName,
		HashedPassword: hashedPassword,
		Token:          tokenString,
		Balance:        1000,
	}
	if err := p.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
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

	if err := tx.Where("username = ?", receiverName).First(&receiver).Error; err != nil {
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


func (p *PostgresRepositoryImpl) UpdateBalance(ctx context.Context, userName string, amount int) error {
	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var purchaser models.User

	if err := tx.Where("userName = ?", userName).First(&purchaser).Error; err != nil {
		tx.Rollback()
		return err
	}

	if purchaser.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient funds for purchase")
	}

	if err := tx.Model(&models.User{}).Where("userName = ?", userName).Update("balance", purchaser.Balance + amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


func (p *PostgresRepositoryImpl) CreatePurchase(ctx context.Context, purchaserName string, merchName string) error {

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

func (p *PostgresRepositoryImpl) FindAppliedTransactions(ctx context.Context, userName string) ([]models.Transaction, error) {
	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var transactions []models.Transaction

	user, err := p.FindUserByName(ctx, userName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("sender_id = ? OR receiver_id = ?", user.ID, user.ID).Find(&transactions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return transactions, tx.Commit().Error
}

func (p *PostgresRepositoryImpl) FindBoughtMerch(ctx context.Context, userName string) ([]responses.InventoryItem, error) {
	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var results []responses.InventoryItem

	user, err := p.FindUserByName(ctx, userName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Table("purchases").
		Select("merches.name AS merch_name, COUNT(*) AS count").
		Joins("JOIN merches ON purchases.merch_id = merches.id").
		Where("purchases.user_id = ?", user.ID).
		Group("merches.name").
		Scan(&results).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return results, tx.Commit().Error
}


func (p *PostgresRepositoryImpl) FindMerchByName(ctx context.Context, merchName string)	 (*models.Merch, error) {
	var merch models.Merch

	if err := p.db.Where("name = ?", merchName).First(&merch).Error; err != nil {
		return nil, err
	}
	return &merch, nil
}
