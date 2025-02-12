package handlers

import (
	"avito/internal/apperror"
	"avito/internal/router/handlers/responses"
	"avito/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)


func BuyHandler(userOperator usecase.UserCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemName := c.Param("item")

		if err := userOperator.BuyItem(itemName); err != nil {
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