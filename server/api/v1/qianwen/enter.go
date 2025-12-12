package qianwen

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	QianwenApi
}

var (
	qianwenService = service.ServiceGroupApp.QianwenServiceGroup.QianwenService
)


