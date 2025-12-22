package localai

type ServiceGroup struct {
	RAGService
}

// 使用单例模式，避免重复创建服务实例
var (
	ragServiceInstance *RAGService
)

func init() {
	ragServiceInstance = NewRAGService()
}

// GetRAGService 获取RAG服务实例
func GetRAGService() *RAGService {
	return ragServiceInstance
}


