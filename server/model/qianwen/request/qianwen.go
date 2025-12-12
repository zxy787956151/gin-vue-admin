package request

// ChatRequest 聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required"` // 用户输入的消息
}


