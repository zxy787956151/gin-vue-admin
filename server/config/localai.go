package config

type LocalAI struct {
	// 本地LLM配置
	LLM LocalLLM `mapstructure:"llm" json:"llm" yaml:"llm"`
	
	// 嵌入模型配置
	Embedding EmbeddingModel `mapstructure:"embedding" json:"embedding" yaml:"embedding"`
	
	// 向量数据库配置
	VectorStore VectorStoreConfig `mapstructure:"vector-store" json:"vector-store" yaml:"vector-store"`
	
	// 训练配置
	Training TrainingConfig `mapstructure:"training" json:"training" yaml:"training"`
	
	// 知识库配置
	Knowledge KnowledgeConfig `mapstructure:"knowledge" json:"knowledge" yaml:"knowledge"`
}

type LocalLLM struct {
	Backend  string `mapstructure:"backend" json:"backend" yaml:"backend"`       // ollama, llama.cpp, vllm
	BaseURL  string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`    // http://localhost:11434
	Model    string `mapstructure:"model" json:"model" yaml:"model"`             // qwen2.5:7b, llama3.1:8b
	Timeout  int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`       // 请求超时时间（秒）
	MaxTokens int   `mapstructure:"max-tokens" json:"max-tokens" yaml:"max-tokens"` // 最大生成token数
}

type EmbeddingModel struct {
	Type    string `mapstructure:"type" json:"type" yaml:"type"`          // local, api
	BaseURL string `mapstructure:"base-url" json:"base-url" yaml:"base-url"` // http://localhost:5001
	Model   string `mapstructure:"model" json:"model" yaml:"model"`       // all-MiniLM-L6-v2
	Dimension int  `mapstructure:"dimension" json:"dimension" yaml:"dimension"` // 向量维度
}

type VectorStoreConfig struct {
	Type       string `mapstructure:"type" json:"type" yaml:"type"`             // chroma, qdrant, local
	Host       string `mapstructure:"host" json:"host" yaml:"host"`             // localhost:8000
	Collection string `mapstructure:"collection" json:"collection" yaml:"collection"` // knowledge_base
	DataPath   string `mapstructure:"data-path" json:"data-path" yaml:"data-path"`   // ./data/vector
}

type TrainingConfig struct {
	Enabled       bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	DataPath      string `mapstructure:"data-path" json:"data-path" yaml:"data-path"`         // ./data/training
	OutputPath    string `mapstructure:"output-path" json:"output-path" yaml:"output-path"`   // ./models/finetuned
	AutoTrain     bool   `mapstructure:"auto-train" json:"auto-train" yaml:"auto-train"`      // 是否自动训练
	MinExamples   int    `mapstructure:"min-examples" json:"min-examples" yaml:"min-examples"` // 触发训练的最小样本数
	Epochs        int    `mapstructure:"epochs" json:"epochs" yaml:"epochs"`
	BatchSize     int    `mapstructure:"batch-size" json:"batch-size" yaml:"batch-size"`
	LearningRate  float64 `mapstructure:"learning-rate" json:"learning-rate" yaml:"learning-rate"`
}

type KnowledgeConfig struct {
	DataPath     string   `mapstructure:"data-path" json:"data-path" yaml:"data-path"`       // ./knowledge_base
	SupportedExt []string `mapstructure:"supported-ext" json:"supported-ext" yaml:"supported-ext"` // [".txt", ".md", ".pdf"]
	ChunkSize    int      `mapstructure:"chunk-size" json:"chunk-size" yaml:"chunk-size"`    // 文本分块大小
	ChunkOverlap int      `mapstructure:"chunk-overlap" json:"chunk-overlap" yaml:"chunk-overlap"` // 分块重叠大小
}

