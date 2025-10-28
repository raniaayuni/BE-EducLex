package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"-"`
	GoogleID string             `bson:"google_id,omitempty" json:"google_id"`
	Role     string             `bson:"role,omitempty" json:"role"`
	Token    string             `bson:"token,omitempty" json:"token"`
}

type Question struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama       string             `bson:"nama" json:"nama"`
	Email      string             `bson:"email" json:"email"`
	Pertanyaan string             `bson:"pertanyaan" json:"pertanyaan"`
	Jawaban    string             `bson:"jawaban,omitempty" json:"jawaban,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type Article struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Image     string             `json:"image" bson:"image"`
	File      string             `json:"file" bson:"file"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Tulisan struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Penulis  string             `bson:"penulis" json:"penulis"`
	Kategori string             `bson:"kategori" json:"kategori"`
	Judul    string             `bson:"judul" json:"judul"`
	Isi      string             `bson:"isi" json:"isi"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Peraturan struct {
    ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
    Judul     string    `bson:"judul" json:"judul"`
    Isi       string    `bson:"isi" json:"isi"`
    Kategori  string    `bson:"kategori" json:"kategori"`
    CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type TokenBlacklist struct {
	Token     string    `bson:"token"`
	ExpiredAt time.Time `bson:"expired_at"`
}