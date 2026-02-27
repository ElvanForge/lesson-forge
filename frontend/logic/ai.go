package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AIProvider interface {
	GenerateContent(ctx context.Context, prompt string, genImage bool) (string, error)
}

// GeminiProvider implementation
type GeminiProvider struct {
	APIKey string
}

func (g *GeminiProvider) GenerateContent(ctx context.Context, prompt string, genImage bool) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", g.APIKey)

	modalities := []string{"TEXT"}
	if genImage {
		modalities = append(modalities, "IMAGE")
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]interface{}{{"text": prompt}}},
		},
		"generationConfig": map[string]interface{}{
			"response_modalities": modalities,
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct { Text string `json:"text"` } `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", fmt.Errorf("AI returned empty content")
}

// FIX: Added missing MockProvider struct definition
type MockProvider struct{}

func (m *MockProvider) GenerateContent(ctx context.Context, p string, img bool) (string, error) {
	return "# Mock Content\nThis is a generated lesson plan for testing purposes.", nil
}

func GetAIProvider(countryCode string) AIProvider {
	if os.Getenv("MOCK_AI") == "true" {
		return &MockProvider{}
	}
	return &GeminiProvider{APIKey: os.Getenv("GEMINI_KEY")}
}