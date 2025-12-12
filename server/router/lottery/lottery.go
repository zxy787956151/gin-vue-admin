package lottery

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type LotteryRouter struct{}

// InitLotteryRouter 初始化彩票路由（公开路由，不需要登录验证）
func (l *LotteryRouter) InitLotteryRouter(Router *gin.RouterGroup) {
	lotteryRouter := Router.Group("lottery")
	lotteryApi := v1.ApiGroupApp.LotteryApiGroup
	{
		lotteryRouter.POST("check", lotteryApi.CheckLottery) // 检查彩票中奖
	}
}

