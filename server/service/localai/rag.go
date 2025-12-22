package localai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/localai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai/response"
)

// RAGService RAG引擎服务
type RAGService struct {
	llm         *LLMService
	vectorStore *VectorStoreService
	training    *TrainingService
}

func NewRAGService() *RAGService {
	return &RAGService{
		llm:         NewLLMService(),
		vectorStore: NewVectorStoreService(),
		training:    NewTrainingService(),
	}
}

// Chat RAG聊天
func (s *RAGService) Chat(ctx context.Context, req request.ChatRequest) (response.ChatResponse, error) {
	startTime := time.Now()

	var answer string
	var sources []localai.SearchResult
	var err error

	if req.UseRAG {
		// 使用 RAG 模式
		answer, sources, err = s.chatWithRAG(ctx, req)
	} else {
		// 直接对话模式
		answer, err = s.chatDirect(ctx, req.Message)
	}

	if err != nil {
		return response.ChatResponse{}, err
	}

	// 异步收集训练数据
	if req.UseRAG && len(sources) > 0 {
		contexts := make([]string, len(sources))
		for i, src := range sources {
			contexts[i] = src.Document.Content
		}
		go s.training.AddExample(
			req.Message,
			strings.Join(contexts, "\n"),
			answer,
			4.0, // 默认评分
			"",
			"auto_collect",
		)
	}

	return response.ChatResponse{
		Answer:    answer,
		Sources:   sources,
		UseRAG:    req.UseRAG,
		SessionID: req.SessionID,
		Time:      time.Since(startTime).String(),
	}, nil
}

// chatWithRAG 使用 RAG 检索增强的对话
func (s *RAGService) chatWithRAG(ctx context.Context, req request.ChatRequest) (string, []localai.SearchResult, error) {
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	// 1. 检索相关文档
	results, err := s.vectorStore.Search(ctx, req.Message, topK)
	if err != nil {
		return "", nil, fmt.Errorf("检索失败: %w", err)
	}

	// 2. 构建增强提示
	contexts := make([]string, len(results))
	for i, result := range results {
		contexts[i] = fmt.Sprintf("[文档%d] %s", i+1, result.Document.Content)
	}

	prompt := s.buildRAGPrompt(req.Message, contexts)

	// 3. 调用 LLM 生成回答
	messages := []localai.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个智能助手，请基于提供的上下文信息回答用户的问题。如果上下文中没有相关信息，请如实告知。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	answer, err := s.llm.Chat(ctx, messages)
	if err != nil {
		return "", nil, fmt.Errorf("生成回答失败: %w", err)
	}

	return answer, results, nil
}

// chatDirect 直接对话
func (s *RAGService) chatDirect(ctx context.Context, message string) (string, error) {
	messages := []localai.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个有帮助的AI助手。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	return s.llm.Chat(ctx, messages)
}

// buildRAGPrompt 构建 RAG 提示
func (s *RAGService) buildRAGPrompt(question string, contexts []string) string {
	var prompt strings.Builder

	prompt.WriteString("基于以下上下文信息回答问题：\n\n")
	prompt.WriteString("【上下文】\n")
	prompt.WriteString(strings.Join(contexts, "\n\n"))
	prompt.WriteString("\n\n【问题】\n")
	prompt.WriteString(question)
	prompt.WriteString("\n\n请根据上下文提供准确、详细的回答：")

	return prompt.String()
}

// IngestDocument 导入文档到知识库
func (s *RAGService) IngestDocument(ctx context.Context, req request.IngestRequest) (response.IngestResponse, error) {
	docID, err := s.vectorStore.AddDocument(ctx, req.Content, req.Metadata)
	if err != nil {
		return response.IngestResponse{}, err
	}

	return response.IngestResponse{
		DocID:     docID,
		Status:    "success",
		Message:   "文档已成功添加到知识库",
		CreatedAt: time.Now(),
	}, nil
}

// SearchDocuments 搜索文档
func (s *RAGService) SearchDocuments(ctx context.Context, req request.SearchRequest) (response.SearchResponse, error) {
	startTime := time.Now()

	topK := req.TopK
	if topK <= 0 {
		topK = 10
	}

	results, err := s.vectorStore.Search(ctx, req.Query, topK)
	if err != nil {
		return response.SearchResponse{}, err
	}

	return response.SearchResponse{
		Results: results,
		Total:   len(results),
		Time:    time.Since(startTime).String(),
	}, nil
}

// SubmitFeedback 提交用户反馈
func (s *RAGService) SubmitFeedback(ctx context.Context, req request.FeedbackRequest) error {
	_, err := s.training.AddExample(
		req.Question,
		"", // 上下文可以从日志中获取
		req.Answer,
		req.Rating,
		req.Feedback,
		"user_feedback",
	)
	return err
}

// StartTraining 开始训练
func (s *RAGService) StartTraining(ctx context.Context, req request.TrainRequest) error {
	return s.training.StartTraining(ctx)
}

// GetTrainingStatus 获取训练状态
func (s *RAGService) GetTrainingStatus() response.TrainStatusResponse {
	isTraining, progress, status := s.training.GetTrainingStatus()

	resp := response.TrainStatusResponse{
		Status:   status,
		Progress: progress,
	}

	if isTraining {
		resp.Message = fmt.Sprintf("训练进行中... %.1f%%", progress)
	} else if status == "completed" {
		resp.Message = "训练已完成"
	} else {
		resp.Message = "未在训练"
	}

	return resp
}

// GetStats 获取系统统计
func (s *RAGService) GetStats() response.StatsResponse {
	vectorStats := s.vectorStore.GetStats()
	trainingStats := s.training.GetStats()

	return response.StatsResponse{
		TotalDocuments: vectorStats["total_documents"].(int),
		TotalExamples:  trainingStats["total_examples"].(int),
		ModelInfo: response.ModelInfo{
			Name:    s.llm.model,
			Backend: s.llm.backend,
			Status:  "running",
		},
		VectorStoreStatus: "healthy",
	}
}

// CheckHealth 健康检查
func (s *RAGService) CheckHealth(ctx context.Context) error {
	return s.llm.CheckHealth(ctx)
}

// DeleteDocument 删除文档
func (s *RAGService) DeleteDocument(ctx context.Context, docID string) error {
	return s.vectorStore.DeleteDocument(ctx, docID)
}


