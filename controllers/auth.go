package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// START login Google
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// contoh login handler
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔹 Cek user di database
	// misalnya hardcode dulu
	email := req.Email
	password := req.Password

	// contoh validasi sederhana
	if email != "admin@example.com" || password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 🔹 Misalnya role ditentukan berdasarkan user
	role := "admin"

	// 🔹 Generate JWT
	jwtToken, err := middleware.GenerateJWT(email, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"email": email,
		"role":  role,
	})
}

// Callback dari Google
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	// Tukar code dengan access token
	tok, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// Ambil data user dari Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var g map[string]any
	_ = json.Unmarshal(body, &g)

	email, _ := g["email"].(string)
	name, _ := g["name"].(string)

	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No email from Google"})
		return
	}

	// ✅ Cek apakah email sudah ada di DB
	var existing models.User
	err = config.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existing)
	if err != nil {
		// Kalau user belum terdaftar → tolak
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not registered",
			"message": "Silakan daftar dulu sebelum login dengan Google",
		})
		return
	}

	// Kalau ada → generate JWT
	jwtToken, err := middleware.GenerateJWT(email, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Login sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"email":   email,
		"name":    name,
		"token":   jwtToken,
	})
}

// END login Google

// Contoh route proteksi
func ProfileHandler(c *gin.Context) {
	email := c.GetString("email")
	name := c.GetString("name")
	c.JSON(http.StatusOK, gin.H{
		"email": email,
		"name":  name,
	})
}
