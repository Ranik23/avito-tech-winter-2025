package handlers

import (
	"avito/internal/apperror"
	"avito/internal/router/handlers/responses"
	"avito/internal/service"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func BuyHandler(manager *service.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		itemName := c.Param("item")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		user, err := manager.AuthService.VerifyToken(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		purchaserName := user.Username

		if err := manager.MerchService.Buy(ctx, purchaserName, itemName); err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, apperror.ErrItemNotFound) {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Успешный ответ"})
	}
}
