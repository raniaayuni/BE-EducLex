package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}

func main() {
	// Konek DB
	config.ConnectDB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert dummy user
	user := TestUser{
		Username: "tester",
		Email:    "tester@mail.com",
	}

	insertResult, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal("❌ Gagal insert:", err)
	}
	fmt.Println("✅ Inserted user with ID:", insertResult.InsertedID)

	// Cek apakah data masuk
	var result TestUser
	err = config.UserCollection.FindOne(ctx, bson.M{"email": "tester@mail.com"}).Decode(&result)
	if err != nil {
		log.Fatal("❌ Gagal find:", err)
	}
	fmt.Println("🎉 User ditemukan di MongoDB:", result.Username, "| Email:", result.Email)
}
