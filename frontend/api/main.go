package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ElvanForge/lesson-forge/backend/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Handler(w http.ResponseWriter, r *http.Request) {
	if pool == nil {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			http.Error(w, "Database connection failed", 500)
			return
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "https://forge.vaelia.app")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api")
	switch {
	case path == "/generate" && r.Method == "POST":
		authMiddleware(http.HandlerFunc(handleGenerate)).ServeHTTP(w, r)
	case path == "/user/credits" && r.Method == "GET":
		authMiddleware(http.HandlerFunc(handleGetCredits)).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 110*time.Second)
	defer cancel()

	userID := r.Header.Get("X-User-ID")
	var req struct {
		Prompt string `json:"prompt"`
		Mode   string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	creditCost := 1
	if req.Mode == "ppt" { creditCost = 2 }

	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2::uuid AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Insufficient credits", 402)
		return
	}

	provider := internal.GetAIProvider()
	content, err := provider.GenerateContent(ctx, req.Prompt)
	if err != nil {
		http.Error(w, "AI Generation Failed", 500)
		return
	}

	var fileData []byte
	var fileName string
	var contentType string

	if req.Mode == "ppt" {
		fileData, fileName, err = internal.GeneratePPTX(userID, content)
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	} else {
		fileData, fileName, err = internal.GeneratePDF(userID, content)
		contentType = "application/pdf"
	}

	if err != nil {
		http.Error(w, "File Generation Failed", 500)
		return
	}

	publicURL, err := uploadToSupabase(fileData, fileName, contentType)
	if err != nil {
		http.Error(w, "Upload Failed", 500)
		return
	}

	pool.Exec(ctx, "INSERT INTO generations (user_id, prompt, file_path, status) VALUES ($1::uuid, $2, $3, 'completed')", userID, req.Prompt, publicURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file": publicURL})
}

func handleGetCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var balance int
	err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"credits": balance})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", 401)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		req, _ := http.NewRequest("GET", os.Getenv("SUPABASE_URL")+"/auth/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			http.Error(w, "Auth failed", 401)
			return
		}
		defer resp.Body.Close()

		var user struct{ ID string `json:"id"` }
		json.NewDecoder(resp.Body).Decode(&user)
		r.Header.Set("X-User-ID", user.ID)
		next.ServeHTTP(w, r)
	})
}

func uploadToSupabase(fileBytes []byte, fileName string, contentType string) (string, error) {
	bucket := "generated-files"
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", os.Getenv("SUPABASE_URL"), bucket, fileName)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(fileBytes))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	if resp.StatusCode != 200 { return "", fmt.Errorf("upload status: %d", resp.StatusCode) }
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucket, fileName), nil
}