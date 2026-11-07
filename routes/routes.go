package routes

import (
	"ledger/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB) {
	h := handlers.TransactionHandler{
		DB: db,
	}

	router.GET("/transactions", h.List)
	router.POST("/transactions", h.Create)
	router.PUT("/transactions/:id", h.Update)
	router.DELETE("/transactions/:id", h.Delete)
}
