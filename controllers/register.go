package controllers

import (
	"net/http"

	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: simpan user ke database MongoDB
	// sementara mock aja
	c.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
		"user":    user,
	})
}
