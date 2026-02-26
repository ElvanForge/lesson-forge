package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	var req struct {
		Prompt string `json:"prompt"`
		Mode   string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", 400)
		return
	}

	// DEBUG CHECK 1: AI Provider
	provider := logic.GetAIProvider()
	if provider == nil {
		log.Println("CRITICAL: GetAIProvider returned nil")
		http.Error(w, "BACKEND_ERROR: AI Provider not initialized. Check OPENAI_API_KEY in Vercel.", 500)
		return
	}

	// DEBUG CHECK 2: AI Generation
	content, err := provider.GenerateContent(r.Context(), req.Prompt)
	if err != nil {
		log.Printf("AI ERROR: %v", err)
		http.Error(w, fmt.Sprintf("AI_ERROR: %v", err), 500)
		return
	}

	var data []byte
	var name string
	// Generate File
	if req.Mode == "ppt" {
		data, name, err = logic.GeneratePPTX(userID, content)
	} else {
		data, name, err = logic.GeneratePDF(userID, content)
	}

	// DEBUG CHECK 3: Document Builder
	if err != nil {
		log.Printf("DOC GEN ERROR: %v", err)
		http.Error(w, "DOC_GEN_ERROR: Logic package failed to build file", 500)
		return
	}

	// DEBUG CHECK 4: Storage
	url, err := uploadToSupabase(data, name, req.Mode)
	if err != nil {
		log.Printf("STORAGE ERROR: %v", err)
		http.Error(w, "STORAGE_ERROR: "+err.Error(), 500)
		return
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

func uploadToSupabase(fileBytes []byte, fileName string, mode string) (string, error) {
	cType := "application/pdf"
	if mode == "ppt" {
		cType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	}

	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if key == "" {
		return "", fmt.Errorf("MISSING_SERVICE_ROLE_KEY")
	}

	url := fmt.Sprintf("%s/storage/v1/object/generated-files/%s", os.Getenv("SUPABASE_URL"), fileName)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(fileBytes))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Content-Type", cType)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Upload failed status: %d", resp.StatusCode)
	}

	return fmt.Sprintf("%s/storage/v1/object/public/generated-files/%s", os.Getenv("SUPABASE_URL"), fileName), nil
}