package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

func main() {
	_ = godotenv.Load("../.env")

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("POST /api/generate", authMiddleware(http.HandlerFunc(handleGenerate)))
	mux.Handle("GET /api/user/credits", authMiddleware(http.HandlerFunc(handleGetCredits)))

	fmt.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

func uploadToSupabase(fileBytes []byte, fileName string, contentType string) (string, error) {
	bucket := "generated-files"
	supabaseURL := os.Getenv("SUPABASE_URL")
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, fileName)

	req, _ := http.NewRequest("POST", uploadURL, bytes.NewReader(fileBytes))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("upload failed: %d", resp.StatusCode)
	}

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucket, fileName), nil
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 110*time.Second)
	defer cancel()

	userID := r.Header.Get("X-User-ID")
	var req struct {
		Prompt string `json:"prompt"`
		Mode   string `json:"mode"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	creditCost := 1
	if req.Mode == "ppt" { creditCost = 2 }

	// Deduct credits and check balance
	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2::uuid AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(map[string]string{"error": "Insufficient credits"})
		return
	}

	provider := GetAIProvider()
	content, _ := provider.GenerateContent(ctx, req.Prompt)

	var fileData []byte
	var fileName string
	var contentType string

	if req.Mode == "ppt" {
		localPath, _ := GeneratePPTX(userID, content)
		fileData, _ = os.ReadFile(localPath)
		fileName = filepath.Base(localPath)
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	} else {
		localPath, _ := GeneratePDF(userID, content)
		fileData, _ = os.ReadFile(localPath)
		fileName = filepath.Base(localPath)
		contentType = "application/pdf"
	}

	publicURL, err := uploadToSupabase(fileData, fileName, contentType)
	if err != nil {
		http.Error(w, "Cloud upload failed", 500)
		return
	}

	pool.Exec(ctx, "INSERT INTO generations (user_id, prompt, file_path, status) VALUES ($1::uuid, $2, $3, 'completed')", userID, req.Prompt, publicURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file": publicURL})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 8 { http.Error(w, "Unauthorized", 401); return }
		token := authHeader[7:]
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", os.Getenv("SUPABASE_URL")+"/auth/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 { http.Error(w, "Auth failed", 401); return }
		defer resp.Body.Close()
		var user struct{ ID string `json:"id"` }
		json.NewDecoder(resp.Body).Decode(&user)
		r.Header.Set("X-User-ID", user.ID)
		next.ServeHTTP(w, r)
	})
}

func handleGetCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var balance int
	pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	json.NewEncoder(w).Encode(map[string]int{"credits": balance})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, apikey")
		if r.Method == "OPTIONS" { w.WriteHeader(http.StatusOK); return }
		next.ServeHTTP(w, r)
	})
}