package models

import "time"

type Purchase struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"column:user_id;index"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MerchID   uint      `gorm:"column:merch_id;index"`
	Merch     Merch     `gorm:"foreignKey:MerchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price     int       `gorm:"column:price;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
