package handlers

import (
	"avito/internal/router/handlers/requests"
	"avito/internal/router/handlers/responses"
	"avito/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)



func AuthHandler(userOperator usecase.UserCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.AuthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}


		token, err := userOperator.Authenticate(req.UserName, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, responses.AuthResponse{
			Token : token,
		})
	}
}