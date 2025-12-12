package qianwen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/qianwen/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/qianwen/response"
)

type QianwenService struct{}

// 通义千问 API 请求结构
type qianwenAPIRequest struct {
	Model string `json:"model"`
	Input struct {
		Messages []Message `json:"messages"`
	} `json:"input"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 通义千问 API 响应结构
type qianwenAPIResponse struct {
	Output struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"output"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
		InputTokens  int `json:"input_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

// Chat 调用通义千问进行聊天
func (q *QianwenService) Chat(req request.ChatRequest) (response.ChatResponse, error) {
	var resp response.ChatResponse

	// 检查配置
	if global.GVA_CONFIG.Qianwen.ApiKey == "" {
		return resp, errors.New("通义千问 API Key 未配置")
	}

	// 如果没有配置模型，使用默认模型
	model := global.GVA_CONFIG.Qianwen.Model
	if model == "" {
		model = "qwen-turbo" // 免费模型
	}

	// 构建请求
	apiReq := qianwenAPIRequest{
		Model: model,
	}
	apiReq.Input.Messages = []Message{
		{
			Role:    "user",
			Content: req.Message,
		},
	}

	// 转换为 JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return resp, fmt.Errorf("构建请求失败: %v", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", bytes.NewBuffer(jsonData))
	if err != nil {
		return resp, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+global.GVA_CONFIG.Qianwen.ApiKey)

	// 发送请求
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return resp, fmt.Errorf("请求失败: %v", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("API 返回错误: %d, %s", httpResp.StatusCode, string(body))
	}

	// 解析响应
	var apiResp qianwenAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return resp, fmt.Errorf("解析响应失败: %v", err)
	}

	// 返回结果
	resp.Reply = apiResp.Output.Text
	return resp, nil
}


