package response

import (
	"time"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai"
)

// ChatResponse RAG聊天响应
type ChatResponse struct {
	Answer    string                   `json:"answer"`              // AI回答
	Sources   []localai.SearchResult   `json:"sources,omitempty"`   // 引用的文档来源
	UseRAG    bool                     `json:"use_rag"`             // 是否使用了RAG
	SessionID string                   `json:"session_id"`          // 会话ID
	Time      string                   `json:"time"`                // 耗时
}

// IngestResponse 知识库导入响应
type IngestResponse struct {
	DocID     string    `json:"doc_id"`      // 文档ID
	Status    string    `json:"status"`      // 状态
	Message   string    `json:"message"`     // 消息
	CreatedAt time.Time `json:"created_at"`  // 创建时间
}

// SearchResponse 向量搜索响应
type SearchResponse struct {
	Results []localai.SearchResult `json:"results"` // 搜索结果
	Total   int                    `json:"total"`   // 总数
	Time    string                 `json:"time"`    // 耗时
}

// TrainStatusResponse 训练状态响应
type TrainStatusResponse struct {
	Status      string    `json:"status"`       // training, completed, failed
	Progress    float32   `json:"progress"`     // 进度 0-100
	Message     string    `json:"message"`      // 消息
	StartTime   time.Time `json:"start_time"`   // 开始时间
	EndTime     *time.Time `json:"end_time"`    // 结束时间
	Metrics     map[string]float64 `json:"metrics,omitempty"` // 训练指标
}

// StatsResponse 系统统计响应
type StatsResponse struct {
	TotalDocuments    int       `json:"total_documents"`     // 文档总数
	TotalExamples     int       `json:"total_examples"`      // 训练样本总数
	ModelInfo         ModelInfo `json:"model_info"`          // 模型信息
	LastTrainingTime  *time.Time `json:"last_training_time"` // 最后训练时间
	VectorStoreStatus string    `json:"vector_store_status"` // 向量库状态
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name      string `json:"name"`       // 模型名称
	Backend   string `json:"backend"`    // 后端类型
	Status    string `json:"status"`     // 状态
	Version   string `json:"version"`    // 版本
}


