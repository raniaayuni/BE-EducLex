package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

var (
	googleClientID     = "778838656131-jfnap1huoa7igvob44b1159gg0e2q99e.apps.googleusercontent.com"
	googleClientSecret = "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"
	redirectURI        = "http://localhost:8080/auth/google/callback"
	jwtSecret          = []byte("SECRET_KEY_KAMU")
)

// -------- Handler 1: Redirect ke Google Login --------
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth"
	params := url.Values{}
	params.Add("client_id", googleClientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("access_type", "offline")
	params.Add("prompt", "select_account")

	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusTemporaryRedirect)
}

// -------- Handler 2: Callback dari Google --------
func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Tukar code dengan token
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokenResp map[string]interface{}
	json.Unmarshal(body, &tokenResp)

	idToken, ok := tokenResp["id_token"].(string)
	if !ok {
		http.Error(w, "No id_token received", http.StatusInternalServerError)
		return
	}

	// Verifikasi id_token
	payload, err := idtoken.Validate(r.Context(), idToken, googleClientID)
	if err != nil {
		http.Error(w, "Invalid ID Token", http.StatusUnauthorized)
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	// Buat JWT sendiri
	claims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	// Kirim JWT ke frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": jwtString,
	})
}

// -------- Middleware JWT --------
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// -------- Endpoint protected --------
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ini data protected, cuma bisa diakses dengan JWT valid"))
}

// -------- Main --------
func main() {
	http.HandleFunc("/auth/google/login", GoogleLoginHandler)
	http.HandleFunc("/auth/google/callback", GoogleCallbackHandler)
	http.Handle("/protected", JWTMiddleware(http.HandlerFunc(ProtectedHandler)))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
