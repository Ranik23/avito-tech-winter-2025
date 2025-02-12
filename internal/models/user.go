package models

import "time"

type User struct {
	ID           	uint      `gorm:"primaryKey;autoIncrement"`
	Username     	string    `gorm:"column:username;unique;not null"`
	HashedPassword	[]byte	  `gorm:"column:hashedpassword;not null"`
	Token			string	  `gorm:"column:token;not null"`
	RegisteredAt 	time.Time `gorm:"column:registered_at;autoCreateTime"`
	Balance      	int       `gorm:"column:balance;not null;default:1000"`
}
