package controllers

import (
	"arithmetic-calculator/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OperationInput struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

func PerformOperation(c *gin.Context) {
	var input OperationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")

	result, err := services.ExecuteOperation(c, username.(string), input.Type, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
