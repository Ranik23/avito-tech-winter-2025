package models

import "time"

type Transaction struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	SenderID   uint      `gorm:"column:sender_id;index"`
	Sender     User      `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ReceiverID uint      `gorm:"column:receiver_id;index"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amount     int       `gorm:"column:amount;not null;check:amount>0"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}
