package models

import "time"

type User struct {
	ID           	uint      `gorm:"primaryKey;autoIncrement"`
	Username     	string    `gorm:"column:username;unique;not null"`
	HashedPassword	[]byte	  `gorm:"column:hashedPassword;not null"`
	Token			string	  `gorm:"column:token;not null"`
	RegisteredAt 	time.Time `gorm:"column:registered_at;autoCreateTime"`
	Balance      	int       `gorm:"column:balance;not null;default:1000"`
}

type Merch struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"column:name;unique;not null"`
	Price int    `gorm:"column:price;not null"`
}

type Transaction struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	SenderID   uint      `gorm:"column:sender_id;index"`
	Sender     User      `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ReceiverID uint      `gorm:"column:receiver_id;index"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amount     int       `gorm:"column:amount;not null;check:amount>0"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Purchase struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"column:user_id;index"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MerchID   uint      `gorm:"column:merch_id;index"`
	Merch     Merch     `gorm:"foreignKey:MerchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price     int       `gorm:"column:price;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
