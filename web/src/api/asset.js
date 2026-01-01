import service from '@/utils/request'

// @Tags Asset
// @Summary 获取个人资产分布数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /asset/distribution [get]
export const getAssetDistribution = () => {
  return service({
    url: '/asset/distribution',
    method: 'get'
  })
}

// @Tags Asset
// @Summary 获取个人资产分布数据2
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /asset/distribution2 [get]
export const getAssetDistribution2 = () => {
  return service({
    url: '/asset/distribution2',
    method: 'get'
  })
}

// @Tags Asset
// @Summary 获取个人资产分布数据3
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /asset/distribution3 [get]
export const getAssetDistribution3 = () => {
  return service({
    url: '/asset/distribution3',
    method: 'get'
  })
}

// @Tags Asset
// @Summary 根据配置项名称获取资产分布详情数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param itemName query string true "配置项名称: 活钱管理/稳健理财/长期投资/保险保障"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /asset/distributionDetail [get]
export const getAssetDistributionByItemName = (itemName) => {
  return service({
    url: '/asset/distributionDetail',
    method: 'get',
    params: {
      itemName
    }
  })
}
