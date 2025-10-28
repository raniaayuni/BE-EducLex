package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// ==== UNTUK USER SAJA ====
// =========================
func GetAllTulisanPublic(c *gin.Context) {
	cursor, err := config.TulisanCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var tulisan []models.Tulisan
	if err := cursor.All(context.Background(), &tulisan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tulisan)
}

// =========================
// ===== UNTUK ADMIN =======
// =========================

// GET semua tulisan (admin juga bisa lihat)
func GetAllTulisan(c *gin.Context) {
	cursor, err := config.TulisanCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var tulisan []models.Tulisan
	if err := cursor.All(context.Background(), &tulisan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tulisan)
}

// CREATE tulisan (admin only)
func CreateTulisan(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa menambahkan tulisan"})
		return
	}

	var input models.Tulisan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	_, err := config.TulisanCollection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil ditambahkan!"})
}

// UPDATE tulisan (admin only)
func UpdateTulisan(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa mengedit tulisan"})
		return
	}

	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var input models.Tulisan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"judul":     input.Judul,
			"isi":       input.Isi,
			"updatedAt": time.Now(),
		},
	}

	_, err := config.TulisanCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil diperbarui"})
}

// DELETE tulisan (admin only)
func DeleteTulisan(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa menghapus tulisan"})
		return
	}

	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := config.TulisanCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil dihapus"})
}
