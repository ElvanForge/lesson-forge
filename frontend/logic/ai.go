package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ElvanForge/lesson-forge/backend/logic"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Handler(w http.ResponseWriter, r *http.Request) {
	if pool == nil {
		p, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Printf("DB INIT ERROR: %v", err)
			http.Error(w, "Database connection failed", 500)
			return
		}
		pool = p
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api")
	if path == "/generate" && r.Method == "POST" {
		authMiddleware(http.HandlerFunc(handleGenerate)).ServeHTTP(w, r)
		return
	}
	if path == "/user/credits" && r.Method == "GET" {
		authMiddleware(http.HandlerFunc(handleGetCredits)).ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	
	// Detect country from Vercel Edge headers to route AI traffic
	countryCode := r.Header.Get("x-vercel-ip-country")

	var req struct {
		Prompt string `json:"prompt"`
		Mode   string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	// Updated to pass countryCode to the AI selector
	provider := logic.GetAIProvider(countryCode)
	if provider == nil {
		http.Error(w, "AI configuration error", 500)
		return
	}

	content, err := provider.GenerateContent(r.Context(), req.Prompt)
	if err != nil {
		log.Printf("AI ERROR (Country: %s): %v", countryCode, err)
		http.Error(w, fmt.Sprintf("AI generation failed: %v", err), 500)
		return
	}

	var data []byte
	var name string
	var cType string

	if req.Mode == "ppt" {
		data, name, err = logic.GeneratePPTX(userID, content)
		cType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	} else {
		data, name, err = logic.GeneratePDF(userID, content)
		cType = "application/pdf"
	}

	if err != nil {
		log.Printf("DOC GEN ERROR: %v", err)
		http.Error(w, "Document generation failed", 500)
		return
	}

	uniqueName := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
	url, err := uploadToSupabase(data, uniqueName, cType)
	if err != nil {
		log.Printf("STORAGE UPLOAD CRITICAL ERROR: %v", err)
		http.Error(w, "File storage failed", 500)
		return
	}

	tx, err := pool.Begin(r.Context())
	if err == nil {
		defer tx.Rollback(r.Context())
		tx.Exec(r.Context(), "UPDATE users SET credit_balance = credit_balance - 1 WHERE id = $1::uuid", userID)
		tx.Exec(r.Context(), "INSERT INTO generations (user_id, prompt, file_path, type) VALUES ($1, $2, $3, $4)", userID, req.Prompt, url, req.Mode)
		tx.Commit(r.Context())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file": url})
}

func handleGetCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var balance int
	err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	if err != nil {
		balance = 0
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
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}
	
	serviceKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("apikey", anonKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Supabase Status %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucket, fileName), nil
}