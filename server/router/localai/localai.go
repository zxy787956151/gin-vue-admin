package localai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type LocalAIRouter struct{}

// InitLocalAIRouter 初始化本地AI路由（公开路由，不需要登录验证）
func (l *LocalAIRouter) InitLocalAIRouter(Router *gin.RouterGroup) {
	localaiRouter := Router.Group("localai")
	localaiApi := v1.ApiGroupApp.LocalAIApiGroup
	{
		// 聊天接口
		localaiRouter.POST("chat", localaiApi.Chat) // RAG聊天
		
		// 知识库管理
		localaiRouter.POST("ingest", localaiApi.IngestDocument)       // 导入文档
		localaiRouter.POST("search", localaiApi.SearchDocuments)      // 搜索文档
		localaiRouter.POST("document/delete", localaiApi.DeleteDocument) // 删除文档
		
		// 训练相关
		localaiRouter.POST("feedback", localaiApi.SubmitFeedback)     // 提交反馈
		localaiRouter.POST("train/start", localaiApi.StartTraining)   // 开始训练
		localaiRouter.GET("train/status", localaiApi.GetTrainingStatus) // 训练状态
		
		// 系统信息
		localaiRouter.GET("stats", localaiApi.GetStats)     // 统计信息
		localaiRouter.GET("health", localaiApi.HealthCheck) // 健康检查
	}
}

