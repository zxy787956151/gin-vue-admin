package localai

import "github.com/flipped-aurora/gin-vue-admin/server/service/localai"

type ApiGroup struct {
	LocalAIApi
}

// getService 获取最新的服务实例（每次调用都获取，确保使用最新配置）
func getService() *localai.RAGService {
	return localai.GetRAGService()
}


