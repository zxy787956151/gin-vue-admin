package asset

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	AssetApi
}

var (
	assetService = service.ServiceGroupApp.AssetServiceGroup.AssetService
)
