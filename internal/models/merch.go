package models


type Merch struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"column:name;unique;not null"`
	Price int    `gorm:"column:price;not null"`
}

type BoughtMerch struct {
	MerchName string
	Count     int
}