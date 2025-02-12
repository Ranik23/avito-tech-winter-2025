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

func AuthHandler(userOperator service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.AuthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		userName, password := req.UserName, req.Password

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		token, err := userOperator.Authenticate(ctx, userName, password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, responses.AuthResponse{
			Token: token,
		})
	}
}
