package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/lottery"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/qianwen"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	LotteryApiGroup lottery.ApiGroup
	QianwenApiGroup qianwen.ApiGroup
}
