package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
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

var pool *pgxpool.Pool

func main() {
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

	// Public routes
	mux.HandleFunc("POST /api/webhook/stripe", handleStripeWebhook)

	// Protected routes (requires auth & credits)
	mux.Handle("POST /api/generate", creditMiddleware(http.HandlerFunc(handleGeneratePPT)))
	mux.Handle("GET /api/user/credits", authMiddleware(http.HandlerFunc(handleGetCredits)))

	// Chain global middlewares
	handler := securityHeadersMiddleware(mux)
	handler = corsMiddleware(config.AllowedOrigins, handler)
	handler = rateLimitMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server starting on port %s", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Global Security Middlewares

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' https://*.supabase.co; frame-ancestors 'none';")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(allowedOrigins string, next http.Handler) http.Handler {
	origins := strings.Split(allowedOrigins, ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isAllowed := false
		for _, o := range origins {
			if o == origin || o == "*" {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

type visitor struct {
	lastSeen time.Time
	count    int
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Simple rate limit: 60 requests per minute
		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
		} else {
			v.count++
		}

		if v.count > 60 {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// Middleware to check user credits before LLM calls
func creditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Header.Get("X-User-ID") // Simple auth for demo
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var balance int
		err := pool.QueryRow(ctx, "SELECT credit_balance FROM users WHERE id = $1", userID).Scan(&balance)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if balance < 1 {
			http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authMiddleware validates the Supabase JWT
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := authHeader[7:] // Remove "Bearer " prefix

		// Verify token with Supabase Auth API
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", os.Getenv("SUPABASE_URL")+"/auth/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY")) // Use anon key for user verification

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			// Don't leak specific error details to the client
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		var sbUser struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&sbUser); err != nil {
			log.Printf("Failed to decode Supabase user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set User ID in header for subsequent handlers
		r.Header.Set("X-User-ID", sbUser.ID)
		next.ServeHTTP(w, r)
	})
}

func handleGetCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var balance int
	err := pool.QueryRow(r.Context(), "SELECT credit_balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"credits": balance})
}

func handleGeneratePPT(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	userID := r.Header.Get("X-User-ID")
	// region := r.Header.Get("X-Region") // Custom header for region detection (Removed as we use LOCATION env var now)

	var reqBody struct {
		Prompt        string `json:"prompt"`
		Mode          string `json:"mode"` // 'lesson' or 'ppt'
		IncludeImages bool   `json:"include_images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Calculate and deduct credit (atomic update)
	creditCost := 1
	if reqBody.Mode == "ppt" && reqBody.IncludeImages {
		creditCost = 2
	}

	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2 AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Credit deduction failed or insufficient balance", http.StatusPaymentRequired)
		return
	}

	// 2. Select AI Provider based on environment configuration
	provider := GetAIProvider()

	// Enrich prompt with mode context
	enrichedPrompt := fmt.Sprintf("Mode: %s. Prompt: %s", reqBody.Mode, reqBody.Prompt)
	if reqBody.Mode == "ppt" && reqBody.IncludeImages {
		enrichedPrompt += " (Include descriptions for 4-6 relevant images)"
	}

	aiContent, err := provider.GenerateContent(ctx, enrichedPrompt)
	if err != nil {
		// Refund credit if AI fails
		log.Printf("AI Generation error for user %s: %v", userID, err)
		pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2", creditCost, userID)
		http.Error(w, "AI service currently unavailable", http.StatusBadGateway)
		return
	}

	// 3. Generate PPTX
	filePath, err := GeneratePPTX(userID, aiContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("PPTX Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Log generation
	_, _ = pool.Exec(ctx, "INSERT INTO generations (user_id, prompt, file_path, status) VALUES ($1, $2, $3, $4)",
		userID, reqBody.Prompt, filePath, "completed")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "PPTX generated successfully",
		"file":    filePath,
	})
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Handle credit top-ups
	w.WriteHeader(http.StatusOK)
}
