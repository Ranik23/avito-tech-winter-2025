package middlewares

import (
	"avito/internal/router/handlers/responses"
	"avito/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)




func AuthMiddleware(userOperator usecase.UserCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
                Errors: "Не авторизован",
            })
            c.Abort()
            return
        }

        // Проверяем действительность токена
        user, err := userOperator.VerifyToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
                Errors: "Неавторизован",
            })
            c.Abort()
            return
        }

        // Сохраняем информацию о пользователе в контексте
        c.Set("user", user)

        c.Next()
    }
}
