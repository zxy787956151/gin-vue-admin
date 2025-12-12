package lottery

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	LotteryApi
}

var (
	lotteryService = service.ServiceGroupApp.LotteryServiceGroup.LotteryService
)

