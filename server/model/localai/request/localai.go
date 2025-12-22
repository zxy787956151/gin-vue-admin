package request

// ChatRequest RAG聊天请求
type ChatRequest struct {
	Message    string `json:"message" binding:"required"` // 用户输入的消息
	UseRAG     bool   `json:"use_rag"`                    // 是否使用RAG检索
	TopK       int    `json:"top_k"`                      // 检索返回的文档数量
	Stream     bool   `json:"stream"`                     // 是否流式输出
	SessionID  string `json:"session_id"`                 // 会话ID
}

// IngestRequest 知识库导入请求
type IngestRequest struct {
	Content  string                 `json:"content" binding:"required"` // 文档内容
	Metadata map[string]interface{} `json:"metadata"`                   // 元数据
	Source   string                 `json:"source"`                     // 来源
}

// SearchRequest 向量搜索请求
type SearchRequest struct {
	Query string `json:"query" binding:"required"` // 查询文本
	TopK  int    `json:"top_k"`                    // 返回结果数量
}

// FeedbackRequest 用户反馈请求
type FeedbackRequest struct {
	QuestionID string  `json:"question_id"`             // 问题ID
	Question   string  `json:"question" binding:"required"` // 问题
	Answer     string  `json:"answer" binding:"required"`   // 答案
	Rating     float32 `json:"rating" binding:"required,min=1,max=5"` // 评分 1-5
	Feedback   string  `json:"feedback"`                // 反馈内容
}

// TrainRequest 训练请求
type TrainRequest struct {
	Method       string  `json:"method"`        // lora, full
	Epochs       int     `json:"epochs"`        // 训练轮数
	BatchSize    int     `json:"batch_size"`    // 批次大小
	LearningRate float64 `json:"learning_rate"` // 学习率
	UseGPU       bool    `json:"use_gpu"`       // 是否使用GPU
}

// DeleteDocRequest 删除文档请求
type DeleteDocRequest struct {
	ID string `json:"id" binding:"required"` // 文档ID
}


