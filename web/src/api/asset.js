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
