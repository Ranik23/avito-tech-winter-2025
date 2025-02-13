package handlers

import (
	"avito/internal/router/handlers/responses"
	"avito/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InfoHandler(userOperator service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		user, err := userOperator.VerifyToken(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		coins := user.Balance

		transactionsSent, err := userOperator.GetSentTransactions(ctx, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return 
		}

		transactionsReceived, err := userOperator.GetReceivedTransactions(ctx, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return 
		}

		merchList, err := userOperator.GetBoughtMerch(ctx, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Errors: err.Error(),
			})
			return
		}

		coinHistory := responses.CoinHistory{
			Received: transactionsReceived,
			Sent: transactionsSent,
		}

		c.JSON(http.StatusOK, responses.InfoResponse{
			Coins: coins,
			Inventory: merchList,
			CoinHistory: coinHistory,
		})
	}
}
