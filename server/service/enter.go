package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/asset"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/localai"
	"github.com/flipped-aurora/gin-vue-admin/server/service/lottery"
	"github.com/flipped-aurora/gin-vue-admin/server/service/qianwen"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	LotteryServiceGroup lottery.ServiceGroup
	QianwenServiceGroup qianwen.ServiceGroup
	LocalAIServiceGroup localai.ServiceGroup
	AssetServiceGroup   asset.ServiceGroup
}
