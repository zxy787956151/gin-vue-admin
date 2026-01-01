package asset

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AssetApi struct{}

// GetAssetDistribution
// @Tags      Asset
// @Summary   获取个人资产分布数据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "获取个人资产分布数据"
// @Router    /asset/distribution [get]
func (a *AssetApi) GetAssetDistribution(c *gin.Context) {
	data, err := assetService.GetAssetDistribution()
	if err != nil {
		global.GVA_LOG.Error("获取个人资产分布数据失败!", zap.Error(err))
		response.FailWithMessage("获取个人资产分布数据失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetAssetDistribution2
// @Tags      Asset
// @Summary   获取个人资产分布数据2
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "获取个人资产分布数据2"
// @Router    /asset/distribution2 [get]
func (a *AssetApi) GetAssetDistribution2(c *gin.Context) {
	data, err := assetService.GetAssetDistribution2()
	if err != nil {
		global.GVA_LOG.Error("获取个人资产分布数据2失败!", zap.Error(err))
		response.FailWithMessage("获取个人资产分布数据2失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetAssetDistribution3
// @Tags      Asset
// @Summary   获取个人资产分布数据3
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "获取个人资产分布数据3"
// @Router    /asset/distribution3 [get]
func (a *AssetApi) GetAssetDistribution3(c *gin.Context) {
	data, err := assetService.GetAssetDistribution3()
	if err != nil {
		global.GVA_LOG.Error("获取个人资产分布数据3失败!", zap.Error(err))
		response.FailWithMessage("获取个人资产分布数据3失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetAssetDistributionByItemName
// @Tags      Asset
// @Summary   根据配置项名称获取资产分布详情数据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     itemName  query  string  true  "配置项名称: 活钱管理/稳健理财/长期投资/保险保障"
// @Success   200       {object}  response.Response{data=interface{},msg=string}  "获取资产分布详情数据"
// @Router    /asset/distributionDetail [get]
func (a *AssetApi) GetAssetDistributionByItemName(c *gin.Context) {
	itemName := c.Query("itemName")
	if itemName == "" {
		response.FailWithMessage("配置项名称不能为空", c)
		return
	}
	
	data, err := assetService.GetAssetDistributionByItemName(itemName)
	if err != nil {
		global.GVA_LOG.Error("获取资产分布详情数据失败!", zap.Error(err))
		response.FailWithMessage("获取资产分布详情数据失败", c)
		return
	}
	response.OkWithData(data, c)
}
