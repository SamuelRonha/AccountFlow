package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	accountHandler *AccountHandler,
	transactionHandler *TransactionHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Welcome
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     "AccountFlow",
			"version": "2.0.0",
			"status":  "running",
			"routes": gin.H{
				"POST /api/v1/accounts":              "create account",
				"GET  /api/v1/accounts/:accountId":   "get account by id",
				"POST /api/v1/transactions":          "create transaction",
			},
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/accounts", accountHandler.Create)
		v1.GET("/accounts/:accountId", accountHandler.GetByID)
		v1.POST("/transactions", transactionHandler.Create)
	}

	return r
}
