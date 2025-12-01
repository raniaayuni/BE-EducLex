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

	// Validasi BidangID dan ambil nama Bidang berdasarkan BidangID
	if input.BidangID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
		return
	}

	// Cari Bidang berdasarkan BidangID
	var bidang models.Bidang
	err := config.BidangCollection.FindOne(context.Background(), bson.M{"_id": input.BidangID}).Decode(&bidang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	// Set BidangNama berdasarkan BidangID
	input.BidangNama = bidang.Nama

	// Menambahkan createdAt dan updatedAt
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	// Insert ke database
	_, err = config.TulisanCollection.InsertOne(context.Background(), input)
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

	// Menangani BidangID dan mendapatkan BidangNama
	if input.BidangID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
		return
	}

	// Cari Bidang berdasarkan BidangID
	var bidang models.Bidang
	err := config.BidangCollection.FindOne(context.Background(), bson.M{"_id": input.BidangID}).Decode(&bidang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	// Set BidangNama berdasarkan BidangID
	input.BidangNama = bidang.Nama

	// Update tulisan
	update := bson.M{
		"$set": bson.M{
			"judul":     input.Judul,
			"isi":       input.Isi,
			"bidang_id": input.BidangID,   // Memperbarui BidangID
			"bidang_nama": input.BidangNama, // Memperbarui BidangNama
			"updatedAt": time.Now(),
		},
	}

	_, err = config.TulisanCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
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

// DOWNLOAD file tulisan
func DownloadFile(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	// Cari tulisan berdasarkan ID
	var tulisan models.Tulisan
	err := config.TulisanCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&tulisan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tulisan tidak ditemukan"})
		return
	}

	// Pastikan file ada
	if tulisan.File == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "File tidak ditemukan"})
		return
	}

	// Kirim file
	c.File(tulisan.File)
}
