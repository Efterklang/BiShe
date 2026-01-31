package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// LLMClient handles interactions with OpenAI-compatible APIs
type LLMClient struct {
	ApiKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

// ChatMessage represents a single message in the chat history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the payload sent to the API
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse represents the response from the API
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// NewLLMClient initializes a new client using environment variables
func NewLLMClient() *LLMClient {
	apiKey := os.Getenv("AI_API_KEY")
	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	// Remove trailing slash from BaseURL to avoid double slashes
	baseURL = strings.TrimSuffix(baseURL, "/")

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	log.Printf("Initializing LLM Client: BaseURL=%s, Model=%s", baseURL, model)

	return &LLMClient{
		ApiKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		Client:  &http.Client{Timeout: 120 * time.Second}, // Longer timeout for analysis generation
	}
}

// GenerateAnalysis sends a prompt to the LLM and returns the generated text
func (c *LLMClient) GenerateAnalysis(prompt string) (string, error) {
	if c.ApiKey == "" {
		return "⚠️ AI API Key 未配置。请在服务器环境变量中设置 AI_API_KEY。", nil
	}

	reqBody := ChatRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "你是一位资深的养生店经营顾问，擅长通过数据分析提供商业洞察。请根据用户提供的经营数据（营收、技师负载、热门项目等），生成一份专业的 Markdown 格式经营分析周报。周报应包含：1. 核心指标解读 2. 存在的问题诊断 3. 具体的改进建议（含项目上下架建议）。语气专业、客观且富有建设性。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Ensure BaseURL doesn't end with slash if we append /chat/completions
	// Assuming BaseURL is like "https://api.openai.com/v1"
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response from AI provider")
}

func (c *LLMClient) GenerateText(systemPrompt, userPrompt string) (string, error) {
	if c.ApiKey == "" {
		return "⚠️ AI API Key 未配置。请在服务器环境变量中设置 AI_API_KEY。", nil
	}

	reqBody := ChatRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response from AI provider")
}
