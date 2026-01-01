package asset

import (
	"github.com/gin-gonic/gin"
)

type AssetRouter struct{}

func (a *AssetRouter) InitAssetRouter(Router *gin.RouterGroup) {
	assetRouterWithoutRecord := Router.Group("asset")
	{
		assetRouterWithoutRecord.GET("distribution", assetApi.GetAssetDistribution) // 获取个人资产分布数据
	}
}
