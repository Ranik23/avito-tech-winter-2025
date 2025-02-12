package handlers

import (
	"avito/internal/controller"
	"avito/internal/router/handlers/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InfoHandler(userOperator *controller.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		_, err := userOperator.VerifyToken(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}
	}
}
