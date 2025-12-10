package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CREATE tulisan (admin only)
func CreateTulisan(c *gin.Context) {
	// Memastikan CORS diterapkan
	c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// (untuk jaga-jaga kalau ada OPTIONS)
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(204)
		return
	}

	// Membuat folder uploads jika belum ada
	os.MkdirAll("uploads", os.ModePerm)

	var input models.Tulisan
	// Mengambil data dari Form Data (bukan JSON)
	input.Penulis = c.PostForm("penulis")
	input.Judul = c.PostForm("judul")
	input.Isi = c.PostForm("isi")

	// Mengambil bidang_id dari Form Data (sebagai string)
	bidangID := c.PostForm("bidang_id")
	if bidangID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak boleh kosong"})
		return
	}

	// Mengonversi bidang_id menjadi primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bidangID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
		return
	}
	input.BidangID = objectID

	// Cari Bidang berdasarkan BidangID untuk mendapatkan Nama Bidang
	var bidang models.Bidang
	err = config.BidangCollection.FindOne(context.Background(), bson.M{"_id": input.BidangID}).Decode(&bidang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	// Menambahkan BidangNama dari Bidang yang ditemukan
	input.BidangNama = bidang.Nama

	// Menangani file gambar
	file, err := c.FormFile("gambar")
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err == nil {
			input.File = path
		}
	}

	// Menangani file dokumen
	dokumen, err := c.FormFile("dokumen")
	if err == nil {
		path := "uploads/" + dokumen.Filename
		if err := c.SaveUploadedFile(dokumen, path); err == nil {
			input.File = path
		}
	}

	// Menambahkan createdAt dan updatedAt
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	// Insert ke dalam collection Tulisan
	_, err = config.TulisanCollection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil ditambahkan!"})
}

// GET semua tulisan
func GetAllTulisan(c *gin.Context) {
	// Pastikan CORS diterapkan
	c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// (untuk jaga-jaga kalau ada OPTIONS)
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(204)
		return
	}

	// Query database Tulisan untuk mendapatkan semua data
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

	// Kirimkan data tulisan ke client
	c.JSON(http.StatusOK, tulisan)
}

func UpdateTulisan(c *gin.Context) {
    // Ambil parameter ID dari URL
    id := c.Param("id")

    // Cek apakah ID valid
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }

    // Ambil data input dari form data
    var input models.Tulisan
    input.Judul = c.PostForm("judul")
    input.Isi = c.PostForm("isi")

    // Mengambil bidang_id dari Form Data (sebagai string)
    bidangID := c.PostForm("bidang_id")
    if bidangID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak boleh kosong"})
        return
    }

    // Mengonversi bidang_id menjadi primitive.ObjectID
    objectID, err := primitive.ObjectIDFromHex(bidangID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
        return
    }
    input.BidangID = objectID

    // Menangani file gambar
    file, err := c.FormFile("gambar")
    if err == nil {
        path := "uploads/" + file.Filename
        if err := c.SaveUploadedFile(file, path); err == nil {
            input.File = path
        }
    }

    // Menangani file dokumen
    dokumen, err := c.FormFile("dokumen")
    if err == nil {
        path := "uploads/" + dokumen.Filename
        if err := c.SaveUploadedFile(dokumen, path); err == nil {
            input.File = path
        }
    }

    // Menangani BidangID dan mendapatkan BidangNama
    if input.BidangID.IsZero() {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
        return
    }

    // Cari Bidang berdasarkan BidangID
    var bidang models.Bidang
    err = config.BidangCollection.FindOne(context.Background(), bson.M{"_id": input.BidangID}).Decode(&bidang)
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
            "bidang_id": input.BidangID,
            "bidang_nama": input.BidangNama,
            "updatedAt": time.Now(),
        },
    }

    // Update data tulisan berdasarkan ID
    result, err := config.TulisanCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tulisan tidak ditemukan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil diperbarui"})
}

// DELETE tulisan
func DeleteTulisan(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := config.TulisanCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tulisan berhasil dihapus"})
}

// Fungsi untuk mengunduh file tulisan berdasarkan ID
func DownloadFile(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cari tulisan berdasarkan ID
	var tulisan models.Tulisan
	err = config.TulisanCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&tulisan)
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
