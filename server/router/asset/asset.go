package asset

import (
	"github.com/gin-gonic/gin"
)

type AssetRouter struct{}

func (a *AssetRouter) InitAssetRouter(Router *gin.RouterGroup) {
	assetRouterWithoutRecord := Router.Group("asset")
	{
		assetRouterWithoutRecord.GET("distribution", assetApi.GetAssetDistribution)           // 获取个人资产分布数据（使用asset配置）
		assetRouterWithoutRecord.GET("distribution2", assetApi.GetAssetDistribution2)         // 获取个人资产分布数据2（使用asset2配置）
		assetRouterWithoutRecord.GET("distribution3", assetApi.GetAssetDistribution3)         // 获取个人资产分布数据3（使用asset3配置）
		assetRouterWithoutRecord.GET("distributionDetail", assetApi.GetAssetDistributionByItemName) // 根据配置项名称获取资产分布详情数据
	}
}
