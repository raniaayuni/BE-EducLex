package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

// --- Redirect ke Google ---
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// --- Callback dari Google ---
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Tukar code dengan token
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Ambil id_token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token received"})
		return
	}

	// Verifikasi id_token (pakai context dari request Gin)
	payload, err := idtoken.Validate(c.Request.Context(), rawIDToken, config.GoogleOauthConfig.ClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID Token"})
		return
	}

	// Ambil data user dari Google
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	// Cari user di DB
	var user models.User
	err = config.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// User baru → register otomatis
		newUser := models.User{
			Email:     email,
			Username:  name,
			Password:  "", // kosong karena login Google
			Provider:  "google",
			CreatedAt: time.Now(),
		}
		_, err := config.UserCollection.InsertOne(context.Background(), newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
			return
		}
		user = newUser
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"email": user.Email,
		"name":  user.Username,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, _ := jwtToken.SignedString(jwtSecret)

	// Login sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Login with Google success",
		"token":   jwtString,
		"user":    user,
	})
}
