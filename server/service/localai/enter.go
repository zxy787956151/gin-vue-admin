package localai

type ServiceGroup struct {
	RAGService
}

// 使用单例模式，避免重复创建服务实例
var (
	ragServiceInstance *RAGService
)

// GetRAGService 获取RAG服务实例（延迟初始化）
func GetRAGService() *RAGService {
	if ragServiceInstance == nil {
		ragServiceInstance = NewRAGService()
	}
	return ragServiceInstance
}

// ReloadService 重新加载服务（在配置加载后调用）
func ReloadService() {
	ragServiceInstance = NewRAGService()
}


