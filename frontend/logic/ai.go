package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AIProvider interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

type GeminiProvider struct {
	APIKey string
}

func (g *GeminiProvider) GenerateContent(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite-preview-02-05:generateContent?key=%s", g.APIKey)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode gemini response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		if result.Error.Message != "" {
			return "", fmt.Errorf("gemini api error: %s", result.Error.Message)
		}
		return "", errors.New("gemini returned no content")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}

type DeepSeekProvider struct {
	APIKey string
}

func (d *DeepSeekProvider) GenerateContent(ctx context.Context, prompt string) (string, error) {
	url := "https://api.deepseek.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.APIKey)

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode deepseek response: %v", err)
	}

	// SAFETY: Check for empty choices before indexing to prevent panic
	if len(result.Choices) == 0 {
		if result.Error.Message != "" {
			return "", fmt.Errorf("deepseek api error: %s", result.Error.Message)
		}
		return "", errors.New("deepseek returned no choices")
	}
	return result.Choices[0].Message.Content, nil
}

type MockProvider struct{}

func (m *MockProvider) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return "# Mock Content\nGenerated for: " + prompt, nil
}

func GetAIProvider(countryCode string) AIProvider {
	if os.Getenv("MOCK_AI") == "true" {
		return &MockProvider{}
	}

	// ROUTING: DeepSeek is ONLY for China (CN)
	if countryCode == "CN" {
		if key := os.Getenv("DEEPSEEK_KEY"); key != "" {
			return &DeepSeekProvider{APIKey: key}
		}
	}

	// FALLBACK: Use Gemini for all other countries
	if key := os.Getenv("GEMINI_KEY"); key != "" {
		return &GeminiProvider{APIKey: key}
	}

	return &MockProvider{}
}