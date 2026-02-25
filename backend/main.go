package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	SupabaseURL        string
	SupabaseAnonKey    string
	Port               string
	AllowedOrigins     string
}

var pool *pgxpool.Pool

func main() {
	_ = godotenv.Load("../.env")

	config := Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey: os.Getenv("SUPABASE_ANON_KEY"),
		Port:            os.Getenv("PORT"),
		AllowedOrigins:  os.Getenv("ALLOWED_ORIGINS"),
	}

	if config.Port == "" { config.Port = "8080" }

	var err error
	pool, err = pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle("POST /api/generate", authMiddleware(http.HandlerFunc(handleGenerate)))
	mux.Handle("GET /api/user/credits", authMiddleware(http.HandlerFunc(handleGetCredits)))

	mux.Handle("/output/", http.StripPrefix("/output/", http.FileServer(http.Dir("./output"))))

	fmt.Printf("🚀 Server running on port %s\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, corsMiddleware(config.AllowedOrigins, mux)))
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 110*time.Second)
	defer cancel()

	userID := r.Header.Get("X-User-ID")
	var req struct {
		Prompt        string `json:"prompt"`
		Mode          string `json:"mode"`
		IncludeImages bool   `json:"includeImages"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", 400)
		return
	}

	// Determine cost: 2 credits for PPT, 1 for Lesson Plan
	creditCost := 1
	if req.Mode == "ppt" {
		creditCost = 2
	}

	// 1. DEDUCT CREDITS
	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2::uuid AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		log.Printf("Credit deduction failed for %s", userID)
		http.Error(w, "Insufficient credits", 402)
		return
	}

	// 2. AI GENERATION
	provider := GetAIProvider()
	promptPrefix := "Lesson plan for: "
	if req.Mode == "ppt" {
		// Instructions for the ## slide split logic in pptx.go
		promptPrefix = "Create a professional presentation outline. Use '##' to separate slides. Format: ## Slide Title\nSlide Content. Prompt: "
	}

	content, err := provider.GenerateContent(ctx, promptPrefix+req.Prompt)
	if err != nil {
		_, _ = pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2::uuid", creditCost, userID)
		http.Error(w, "AI Error", 502)
		return
	}

	// 3. FILE GENERATION (Switching based on mode)
	var filePath string
	if req.Mode == "ppt" {
		filePath, err = GeneratePPTX(userID, content) // Calls your new pptx.go function
	} else {
		filePath, err = GeneratePDF(userID, content)
	}

	if err != nil {
		log.Printf("❌ GENERATION ERROR: %v", err)
		_, _ = pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2::uuid", creditCost, userID)
		http.Error(w, "File Generation failed", 500)
		return
	}

	// 4. SAVE TO DATABASE (Persistence)
	_, err = pool.Exec(ctx, `
		INSERT INTO generations (user_id, prompt, file_path, status) 
		VALUES ($1::uuid, $2, $3, $4)`,
		userID, req.Prompt, filePath, "completed",
	)
	if err != nil {
		log.Printf("❌ DB LOG ERROR: %v", err)
	}

	log.Printf("✅ %s Created and Saved: %s", req.Mode, filePath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "completed",
		"file":   "http://localhost:8080/" + filePath,
	})
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
	err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", 404)
			return
		}
		http.Error(w, "DB error", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"credits": balance})
}

func corsMiddleware(origins string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" { w.WriteHeader(200); return }
		next.ServeHTTP(w, r)
	})
}