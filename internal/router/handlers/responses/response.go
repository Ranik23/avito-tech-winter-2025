package responses


type ErrorResponse struct {
	Errors string 			`json:"errors"`
}

type AuthResponse struct {
	Token string			`json:"token"`
}

type InfoResponse struct {
	Coins       int              `json:"coins"`
	Inventory   []InventoryItem  `json:"inventory"`
	CoinHistory CoinHistory      `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type" gorm:"column:merch_name"`
	Quantity int    `json:"quantity" gorm:"column:count"`
}


type CoinHistory struct {
	Received []Transaction 		 `json:"received"`
	Sent     []Transaction 		 `json:"sent"`
}

type Transaction struct {
	FromUser string 			 `json:"fromUser,omitempty"`
	ToUser   string 			 `json:"toUser,omitempty"`
	Amount   int    			 `json:"amount"`
}