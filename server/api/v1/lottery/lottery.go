package lottery

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LotteryApi struct{}

// CheckLottery 检查彩票中奖
// @Tags Lottery
// @Summary 检查彩票中奖
// @Description 批量检查双色球或大乐透是否中奖，返回中奖等级
// @Accept application/json
// @Produce application/json
// @Param data body request.CheckLotteryReq true "彩票类型和号码"
// @Success 200 {object} response.Response{data=response.CheckLotteryResp,msg=string} "检查成功"
// @Router /lottery/check [post]
func (l *LotteryApi) CheckLottery(c *gin.Context) {
	var req request.CheckLotteryReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := lotteryService.CheckLottery(req)
	if err != nil {
		global.GVA_LOG.Error("检查彩票失败!", zap.Error(err))
		response.FailWithMessage("检查彩票失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "检查成功", c)
}
