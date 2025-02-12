package requests



type AuthRequest struct {
	UserName string			`json:"username" binding:"required"`
	Password string			`json:"password" binding:"required"`
}


type SendCoinRequest struct {
	ToUser string			`json:"toUser" binding:"required"`
	Amount int				`json:"amount" binding:"required"`
}