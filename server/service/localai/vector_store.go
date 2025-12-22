package localai

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai"
	"github.com/google/uuid"
)

// VectorStoreService 向量存储服务
type VectorStoreService struct {
	mu        sync.RWMutex
	documents map[string]localai.Document // 内存存储
	dataPath  string
}

func NewVectorStoreService() *VectorStoreService {
	dataPath := global.GVA_CONFIG.LocalAI.VectorStore.DataPath
	if dataPath == "" {
		dataPath = "./data/vector"
	}

	// 创建数据目录
	os.MkdirAll(dataPath, 0755)

	service := &VectorStoreService{
		documents: make(map[string]localai.Document),
		dataPath:  dataPath,
	}

	// 加载已有数据
	service.loadFromDisk()

	return service
}

// AddDocument 添加文档
func (s *VectorStoreService) AddDocument(ctx context.Context, content string, metadata map[string]interface{}) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 生成文档ID
	docID := uuid.New().String()

	// 生成向量（这里简化处理，实际应该调用嵌入模型）
	vector := s.generateSimpleVector(content)

	doc := localai.Document{
		ID:        docID,
		Content:   content,
		Metadata:  metadata,
		Vector:    vector,
		CreatedAt: time.Now(),
	}

	s.documents[docID] = doc

	// 持久化到磁盘
	go s.saveToDisk()

	return docID, nil
}

// Search 向量搜索
func (s *VectorStoreService) Search(ctx context.Context, query string, topK int) ([]localai.SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if topK <= 0 {
		topK = 5
	}

	// 生成查询向量
	queryVector := s.generateSimpleVector(query)

	// 计算相似度并排序
	results := make([]localai.SearchResult, 0)
	for _, doc := range s.documents {
		similarity := s.cosineSimilarity(queryVector, doc.Vector)
		results = append(results, localai.SearchResult{
			Document: doc,
			Score:    similarity,
			Distance: 1 - similarity,
		})
	}

	// 简单排序（实际应该使用更高效的算法）
	s.sortByScore(results)

	// 返回 topK 结果
	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

// DeleteDocument 删除文档
func (s *VectorStoreService) DeleteDocument(ctx context.Context, docID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.documents[docID]; !exists {
		return errors.New("文档不存在")
	}

	delete(s.documents, docID)

	// 持久化
	go s.saveToDisk()

	return nil
}

// GetStats 获取统计信息
func (s *VectorStoreService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_documents": len(s.documents),
		"data_path":       s.dataPath,
	}
}

// generateSimpleVector 生成简单向量（示例实现）
// 实际应该调用嵌入模型服务
func (s *VectorStoreService) generateSimpleVector(text string) []float32 {
	// 这是一个非常简单的向量化实现，仅用于演示
	// 实际应该调用 sentence-transformers 或其他嵌入模型
	dimension := 384 // all-MiniLM-L6-v2 的维度
	vector := make([]float32, dimension)
	
	// 简单的字符频率向量化
	for i, char := range text {
		if i >= dimension {
			break
		}
		vector[i%dimension] += float32(char) / 1000.0
	}

	// 归一化
	var norm float32
	for _, v := range vector {
		norm += v * v
	}
	if norm > 0 {
		norm = float32(1.0 / (norm + 0.0001))
		for i := range vector {
			vector[i] *= norm
		}
	}

	return vector
}

// cosineSimilarity 计算余弦相似度
func (s *VectorStoreService) cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(1.0) * normA * normB)
}

// sortByScore 按分数排序
func (s *VectorStoreService) sortByScore(results []localai.SearchResult) {
	// 简单的冒泡排序
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
}

// saveToDisk 保存到磁盘
func (s *VectorStoreService) saveToDisk() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.dataPath, "documents.json")
	
	data, err := json.MarshalIndent(s.documents, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// loadFromDisk 从磁盘加载
func (s *VectorStoreService) loadFromDisk() error {
	filePath := filepath.Join(s.dataPath, "documents.json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，这是正常的
		}
		return err
	}

	return json.Unmarshal(data, &s.documents)
}

// Clear 清空所有文档
func (s *VectorStoreService) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.documents = make(map[string]localai.Document)
	return s.saveToDisk()
}


