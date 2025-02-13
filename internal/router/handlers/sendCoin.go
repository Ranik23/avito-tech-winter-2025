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

func SendCoinHandler(userOperator service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.SendCoinRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}

		tokenString := c.GetHeader("Authorization")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		user, err := userOperator.VerifyToken(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		userName := user.Username

		if err := userOperator.SendCoins(ctx, userName, req.ToUser, req.Amount); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Успешный ответ"})
	}
}
