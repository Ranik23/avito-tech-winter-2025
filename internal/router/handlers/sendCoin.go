package handlers

import (
	"avito/internal/router/handlers/requests"
	"avito/internal/router/handlers/responses"
	"avito/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SendCoinHandler(userOperator service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.SendCoinRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := userOperator.SendCoins(ctx, req.ToUser, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Успешный ответ"})
	}
}
