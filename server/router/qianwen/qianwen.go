package qianwen

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type QianwenRouter struct{}

// InitQianwenRouter 初始化通义千问路由（公开路由，不需要登录验证）
func (q *QianwenRouter) InitQianwenRouter(Router *gin.RouterGroup) {
	qianwenRouter := Router.Group("qianwen")
	qianwenApi := v1.ApiGroupApp.QianwenApiGroup
	{
		qianwenRouter.POST("chat", qianwenApi.Chat) // 与通义千问聊天
	}
}



