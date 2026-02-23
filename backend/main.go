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
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
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
func handleGenerate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Engineering Guard: Default to lesson mode if not specified
	if reqBody.Mode == "" {
		reqBody.Mode = "lesson"
	}

	creditCost := 1
	if reqBody.IncludeImages {
		creditCost = 2
	}

	// 1. Pre-check balance (Optional but good for UX)
	// 2. Call AI Provider (Mock until Feb 27)
	provider := GetAIProvider()
	aiContent, err := provider.GenerateContent(ctx, reqBody.Prompt)
	if err != nil {
		log.Printf("AI Provider error: %v", err)
		http.Error(w, "AI service unavailable", http.StatusBadGateway)
		return
	}

	// 3. POST-PAY: Deduct credits
	tag, err := pool.Exec(ctx, `
		UPDATE users 
		SET credit_balance = credit_balance - $1 
		WHERE id = $2 AND credit_balance >= $1`, 
		creditCost, userID)
	
	if err != nil || tag.RowsAffected() == 0 {
		http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
		return
	}

	// 4. File Generation Routing
	var filePath string
	var genErr error

	if reqBody.Mode == "lesson" {
		filePath, genErr = GeneratePDF(userID, aiContent)
	} else {
		filePath, genErr = GeneratePPTX(userID, aiContent)
	}

	if genErr != nil {
		log.Printf("Generation failed: %v", genErr)
		// REFUND: Atomic increment
		_, _ = pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2", creditCost, userID)
		http.Error(w, "Failed to build file", http.StatusInternalServerError)
		return
	}

	// 5. Finalize Generation Record
	_, err = pool.Exec(ctx, `
		INSERT INTO generations (user_id, prompt, file_path, status) 
		VALUES ($1, $2, $3, 'completed')`,
		userID, reqBody.Prompt, filePath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"file_path": filePath,
		"mode":      reqBody.Mode,
	})
}
var pool *pgxpool.Pool

func main() {
	// Load environment variables from the root .env file
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

	// Public routes
	mux.HandleFunc("POST /api/webhook/stripe", handleStripeWebhook)

	// Protected routes (requires auth & credits)
	mux.Handle("POST /api/generate", authMiddleware(creditMiddleware(http.HandlerFunc(handleGeneratePPT))))
	mux.Handle("GET /api/user/credits", authMiddleware(http.HandlerFunc(handleGetCredits)))

	// Chain global middlewares
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
	fmt.Printf("Allowed Origins: %s\n", config.AllowedOrigins)
	fmt.Printf("Database URL present: %v\n", config.DatabaseURL != "")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %v\n", err)
		log.Fatal(err)
	}
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
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
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isAllowed := false

		// Always allow if no origin (non-browser)
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

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
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			fmt.Printf("CORS rejected origin: [%s]. Allowed: %v\n", origin, origins)
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
        userID := r.Header.Get("X-User-ID") 
        
        var balance int
        ctx := r.Context()
        err := pool.QueryRow(ctx, "SELECT credit_balance FROM users WHERE id = $1", userID).Scan(&balance)
        
        if err != nil {
            if err == pgx.ErrNoRows {
                fmt.Printf("DEBUG: Database could not find ID: [%s]\n", userID) // Add this log
                http.Error(w, "User not found", http.StatusNotFound)
                log.Printf("DB Error: %v", err)
                return
            }
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        } // <--- You were missing this brace

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
		if err != nil {
			fmt.Printf("Auth middleware network error: %v\n", err)
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Auth middleware Supabase rejection: status %d for token starting with %s\n", resp.StatusCode, token[:10])
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

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

	// 1. Calculate credit cost
	creditCost := 1
	if reqBody.Mode == "ppt" && reqBody.IncludeImages {
		creditCost = 2
	}

	// 2. Double check balance before LLM call (extra safety beyond middleware)
	var balance int
	err := pool.QueryRow(ctx, "SELECT credit_balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil || balance < creditCost {
		http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
		return
	}

	// 3. Select AI Provider and Generate
	provider := GetAIProvider()
	enrichedPrompt := fmt.Sprintf("Mode: %s. Prompt: %s", reqBody.Mode, reqBody.Prompt)
	if reqBody.Mode == "ppt" && reqBody.IncludeImages {
		enrichedPrompt += " (Include descriptions for 4-6 relevant images)"
	}

	aiContent, err := provider.GenerateContent(ctx, enrichedPrompt)
	if err != nil {
		log.Printf("AI Generation error for user %s: %v", userID, err)
		http.Error(w, "AI service currently unavailable", http.StatusBadGateway)
		return
	}

	// 4. Deduct credit ONLY after successful AI response (Post-pay)
	res, err := pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance - $1 WHERE id = $2 AND credit_balance >= $1", creditCost, userID)
	if err != nil || res.RowsAffected() == 0 {
		// This should theoretically not happen if balance was checked above and no concurrent reqs
		http.Error(w, "Credit deduction failed", http.StatusInternalServerError)
		return
	}

	// 5. Generate PPTX
	filePath, err := GeneratePPTX(userID, aiContent)
	if err != nil {
		// Refund if file gen fails? User said "ONLY after a successful 200 OK response from the LLM"
		// If LLM was successful, we deducted. If file fails, we might still want to refund.
		log.Printf("PPTX Generation failed for user %s: %v", userID, err)
		pool.Exec(ctx, "UPDATE users SET credit_balance = credit_balance + $1 WHERE id = $2", creditCost, userID)
		http.Error(w, "Generation failed", http.StatusInternalServerError)
		return
	}

	// 6. Log generation
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
