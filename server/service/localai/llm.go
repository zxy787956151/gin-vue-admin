package localai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai"
)

// LLMService 本地大模型服务
type LLMService struct {
	client  *http.Client
	baseURL string
	model   string
	backend string
}

func NewLLMService() *LLMService {
	config := global.GVA_CONFIG.LocalAI.LLM
	
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second
	}

	return &LLMService{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: config.BaseURL,
		model:   config.Model,
		backend: config.Backend,
	}
}

// Chat 对话接口
func (s *LLMService) Chat(ctx context.Context, messages []localai.ChatMessage) (string, error) {
	if s.baseURL == "" {
		return "", errors.New("本地模型未配置，请在 config.yaml 中配置 local-ai.llm")
	}

	switch s.backend {
	case "ollama":
		return s.chatOllama(ctx, messages)
	case "llama.cpp":
		return s.chatLlamaCpp(ctx, messages)
	case "vllm":
		return s.chatVLLM(ctx, messages)
	default:
		return s.chatOllama(ctx, messages)
	}
}

// chatOllama 使用 Ollama API
func (s *LLMService) chatOllama(ctx context.Context, messages []localai.ChatMessage) (string, error) {
	// 构建 Ollama 格式的请求
	reqBody := map[string]interface{}{
		"model":  s.model,
		"messages": messages,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"num_predict": global.GVA_CONFIG.LocalAI.LLM.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/chat", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API返回错误: %d, %s", resp.StatusCode, string(body))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Message.Content, nil
}

// chatLlamaCpp 使用 llama.cpp API
func (s *LLMService) chatLlamaCpp(ctx context.Context, messages []localai.ChatMessage) (string, error) {
	// 将消息转换为 prompt
	var prompt strings.Builder
	for _, msg := range messages {
		if msg.Role == "system" {
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		} else if msg.Role == "user" {
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		} else if msg.Role == "assistant" {
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
		}
	}
	prompt.WriteString("Assistant: ")

	reqBody := map[string]interface{}{
		"prompt":      prompt.String(),
		"n_predict":   global.GVA_CONFIG.LocalAI.LLM.MaxTokens,
		"temperature": 0.7,
		"stop":        []string{"\nUser:", "\nSystem:"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/completion", bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Content, nil
}

// chatVLLM 使用 vLLM API
func (s *LLMService) chatVLLM(ctx context.Context, messages []localai.ChatMessage) (string, error) {
	// vLLM 兼容 OpenAI API 格式
	reqBody := map[string]interface{}{
		"model":    s.model,
		"messages": messages,
		"max_tokens": global.GVA_CONFIG.LocalAI.LLM.MaxTokens,
		"temperature": 0.7,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
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
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", errors.New("没有返回内容")
}

// CheckHealth 检查模型健康状态
func (s *LLMService) CheckHealth(ctx context.Context) error {
	if s.baseURL == "" {
		return errors.New("模型未配置")
	}

	var endpoint string
	switch s.backend {
	case "ollama":
		endpoint = "/api/tags"
	case "llama.cpp":
		endpoint = "/health"
	case "vllm":
		endpoint = "/v1/models"
	default:
		endpoint = "/api/tags"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("无法连接到模型服务: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("模型服务返回错误: %d", resp.StatusCode)
	}

	return nil
}


