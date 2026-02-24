package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	DatabaseURL        string
	StripeKey          string
	DeepSeekKey        string
	GeminiKey          string
	SupabaseURL        string
	SupabaseServiceKey string
	AllowedOrigins     string
	Port               string
}

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	CreditBalance int    `json:"credit_balance"`
}

type GenerateRequest struct {
	Prompt        string `json:"prompt"`
	Mode          string `json:"mode"`
	IncludeImages bool   `json:"includeImages"`
}

var pool *pgxpool.Pool

func main() {
	_ = godotenv.Load("../.env")

	config := Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		StripeKey:          os.Getenv("STRIPE_KEY"),
		DeepSeekKey:        os.Getenv("DEEPSEEK_KEY"),
		GeminiKey:          os.Getenv("GEMINI_KEY"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		AllowedOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		Port:               os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	var err error
	pool, err = pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/webhook/stripe", handleStripeWebhook)
	mux.Handle("POST /api/generate", authMiddleware(creditMiddleware(http.HandlerFunc(handleGeneratePPT))))
	mux.Handle("GET /api/user/credits", authMiddleware(http.HandlerFunc(handleGetCredits)))

	handler := requestLoggerMiddleware(mux)
	handler = securityHeadersMiddleware(handler)
	handler = corsMiddleware(config.AllowedOrigins, handler)
	handler = rateLimitMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server starting on port %s\n", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// --- HANDLERS ---

func handleGeneratePPT(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	userID := r.Header.Get("X-User-ID")

	var reqBody struct {
		Prompt        string `json:"prompt"`
		Mode          string `json:"mode"`
		IncludeImages bool   `json:"include_images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	creditCost := 1
	if reqBody.Mode == "ppt" && reqBody.IncludeImages {
		creditCost = 2
	}

	// 1. Double check balance
	var balance int
	err := pool.QueryRow(ctx, "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	if err != nil || balance < creditCost {
		http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
		return
	}

	// 2. Generate Content
	provider := GetAIProvider()
	enrichedPrompt := fmt.Sprintf("Mode: %s. Prompt: %s", reqBody.Mode, reqBody.Prompt)
	aiContent, err := provider.GenerateContent(ctx, enrichedPrompt)
	if err != nil {
		log.Printf("AI Generation error: %v", err)
		http.Error(w, "AI service unavailable", http.StatusBadGateway)
		return
	}

	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2::uuid AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Credit deduction failed", http.StatusInternalServerError)
		return
	}

	filePath, err := GeneratePPTX(userID, aiContent)
	if err != nil {
		log.Printf("File Generation failed: %v", err)
		_, _ = pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2::uuid", creditCost, userID)
		http.Error(w, "Generation failed", http.StatusInternalServerError)
		return
	}

	_, err = pool.Exec(ctx, "INSERT INTO generations (user_id, prompt, file_path, status) VALUES ($1::uuid, $2, $3, 'completed')",
		userID, reqBody.Prompt, filePath)
	if err != nil {
		log.Printf("Failed to log generation: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success",
		"file":    filePath,
	})
}

func handleGetCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var balance int
	err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"credits": balance})
}

// --- MIDDLEWARES ---

func creditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		
		// TACTICAL LOGGING: See exactly what ID is hitting the DB
		log.Printf("CREDIT_CHECK: UserID [%s]", userID)

		var balance int
		err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1::uuid", userID).Scan(&balance)
		
		if err != nil {
			log.Printf("DB_ERROR in creditMiddleware: %v", err) // This will tell us the EXACT error
			if err == pgx.ErrNoRows {
				http.Error(w, "User not found in local database", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}	
		if balance < 1 {
			http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := authHeader[7:]

		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", os.Getenv("SUPABASE_URL")+"/auth/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		var sbUser struct {
			ID string `json:"id"`
		}
		json.NewDecoder(resp.Body).Decode(&sbUser)
		r.Header.Set("X-User-ID", sbUser.ID)
		next.ServeHTTP(w, r)
	})
}

// --- STUBS & UTILS ---

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }
func requestLoggerMiddleware(n http.Handler) http.Handler      { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.ServeHTTP(w, r) }) }
func securityHeadersMiddleware(n http.Handler) http.Handler    { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.ServeHTTP(w, r) }) }
func corsMiddleware(allowedOrigins string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hard-code the origin for local development to kill the CORS error
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle the Preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func rateLimitMiddleware(n http.Handler) http.Handler        { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.ServeHTTP(w, r) }) }
