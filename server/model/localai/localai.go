package localai

import "time"

// Document 文档结构
type Document struct {
	ID       string    `json:"id"`
	Content  string    `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
	Vector   []float32 `json:"vector,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Document   Document `json:"document"`
	Score      float32  `json:"score"`
	Distance   float32  `json:"distance"`
}

// TrainingExample 训练样本
type TrainingExample struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Context   string    `json:"context"`
	Answer    string    `json:"answer"`
	Rating    float32   `json:"rating"` // 1-5分
	Feedback  string    `json:"feedback"`
	Source    string    `json:"source"` // user_feedback, auto_collect
	CreatedAt time.Time `json:"created_at"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}


