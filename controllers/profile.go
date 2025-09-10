package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	// ambil data user dari JWT (yang diset di middleware)
	userEmail, _ := c.Get("email")
	userName, _ := c.Get("name")

	c.JSON(http.StatusOK, gin.H{
		"email": userEmail,
		"name":  userName,
	})
}
