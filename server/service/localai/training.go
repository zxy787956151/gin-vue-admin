package localai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TrainingService 训练服务
type TrainingService struct {
	mu           sync.RWMutex
	examples     map[string]localai.TrainingExample
	dataPath     string
	isTraining   bool
	trainingProgress float32
	trainingStatus   string
}

func NewTrainingService() *TrainingService {
	dataPath := global.GVA_CONFIG.LocalAI.Training.DataPath
	if dataPath == "" {
		dataPath = "./data/training"
	}

	// 创建数据目录
	os.MkdirAll(dataPath, 0755)

	service := &TrainingService{
		examples:   make(map[string]localai.TrainingExample),
		dataPath:   dataPath,
		trainingStatus: "idle",
	}

	// 加载已有数据
	service.loadFromDisk()

	return service
}

// AddExample 添加训练样本
func (s *TrainingService) AddExample(question, contextStr, answer string, rating float32, feedback string, source string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	exampleID := uuid.New().String()

	example := localai.TrainingExample{
		ID:        exampleID,
		Question:  question,
		Context:   contextStr,
		Answer:    answer,
		Rating:    rating,
		Feedback:  feedback,
		Source:    source,
		CreatedAt: time.Now(),
	}

	s.examples[exampleID] = example

	// 持久化
	go s.saveToDisk()

	// 检查是否需要触发自动训练
	if global.GVA_CONFIG.LocalAI.Training.AutoTrain && s.shouldTriggerTraining() {
		go s.StartTraining(context.Background())
	}

	return exampleID, nil
}

// GetStats 获取统计信息
func (s *TrainingService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_examples":     len(s.examples),
		"is_training":        s.isTraining,
		"training_progress":  s.trainingProgress,
		"training_status":    s.trainingStatus,
	}
}

// StartTraining 开始训练
func (s *TrainingService) StartTraining(ctx context.Context) error {
	s.mu.Lock()
	if s.isTraining {
		s.mu.Unlock()
		return fmt.Errorf("训练正在进行中")
	}
	s.isTraining = true
	s.trainingProgress = 0
	s.trainingStatus = "training"
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isTraining = false
		s.mu.Unlock()
	}()

	global.GVA_LOG.Info("开始训练模型...", zap.Int("examples", len(s.examples)))

	// 1. 准备训练数据
	if err := s.prepareTrainingData(); err != nil {
		s.mu.Lock()
		s.trainingStatus = "failed"
		s.mu.Unlock()
		return fmt.Errorf("准备训练数据失败: %w", err)
	}

	s.mu.Lock()
	s.trainingProgress = 30
	s.mu.Unlock()

	// 2. 实际训练过程（这里简化处理）
	// 在实际应用中，这里应该调用训练脚本或API
	global.GVA_LOG.Info("训练中...")
	time.Sleep(5 * time.Second) // 模拟训练过程

	s.mu.Lock()
	s.trainingProgress = 80
	s.mu.Unlock()

	// 3. 保存模型
	global.GVA_LOG.Info("保存模型...")
	time.Sleep(2 * time.Second)

	s.mu.Lock()
	s.trainingProgress = 100
	s.trainingStatus = "completed"
	s.mu.Unlock()

	global.GVA_LOG.Info("训练完成！")

	return nil
}

// shouldTriggerTraining 是否应该触发训练
func (s *TrainingService) shouldTriggerTraining() bool {
	minExamples := global.GVA_CONFIG.LocalAI.Training.MinExamples
	if minExamples == 0 {
		minExamples = 100
	}

	return len(s.examples) >= minExamples
}

// prepareTrainingData 准备训练数据
func (s *TrainingService) prepareTrainingData() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 转换为 Alpaca 格式
	type AlpacaExample struct {
		Instruction string `json:"instruction"`
		Input       string `json:"input"`
		Output      string `json:"output"`
	}

	examples := make([]AlpacaExample, 0, len(s.examples))
	for _, ex := range s.examples {
		// 只使用评分较高的样本
		if ex.Rating < 3.0 {
			continue
		}

		alpaca := AlpacaExample{
			Instruction: "根据以下上下文回答问题",
			Input:       fmt.Sprintf("上下文：%s\n\n问题：%s", ex.Context, ex.Question),
			Output:      ex.Answer,
		}
		examples = append(examples, alpaca)
	}

	// 保存为 JSON 文件
	outputPath := filepath.Join(s.dataPath, "train.json")
	data, err := json.MarshalIndent(examples, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// saveToDisk 保存到磁盘
func (s *TrainingService) saveToDisk() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.dataPath, "examples.json")
	
	data, err := json.MarshalIndent(s.examples, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// loadFromDisk 从磁盘加载
func (s *TrainingService) loadFromDisk() error {
	filePath := filepath.Join(s.dataPath, "examples.json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.examples)
}

// GetTrainingStatus 获取训练状态
func (s *TrainingService) GetTrainingStatus() (bool, float32, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.isTraining, s.trainingProgress, s.trainingStatus
}


