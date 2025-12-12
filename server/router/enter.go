package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/localai"
	"github.com/flipped-aurora/gin-vue-admin/server/router/lottery"
	"github.com/flipped-aurora/gin-vue-admin/server/router/qianwen"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Lottery lottery.RouterGroup
	Qianwen qianwen.RouterGroup
	LocalAI localai.RouterGroup
}
