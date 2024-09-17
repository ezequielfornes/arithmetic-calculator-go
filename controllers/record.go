package controllers

import (
	"arithmetic-calculator/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetRecords(c *gin.Context) {
	username, _ := c.Get("username")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "30"))

	paginatedRecords, err := services.FetchRecords(username.(string), page, pageSize)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, paginatedRecords)
}
