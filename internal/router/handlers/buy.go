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

func BuyHandler(userOperator service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		_ = c.GetHeader("Authorization")

		itemName := c.Param("item")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := userOperator.BuyItem(ctx, itemName); err != nil {
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
