package handlers

import (
	"avito/internal/router/handlers/requests"
	"avito/internal/router/handlers/responses"
	"avito/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)



func SendCoinHandler(userOperator usecase.UserCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.SendCoinRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}

		err := userOperator.SendCoins(req.ToUser, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Успешный ответ"})
	}
}