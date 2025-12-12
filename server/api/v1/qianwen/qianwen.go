package qianwen

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/qianwen/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QianwenApi struct{}

// Chat 与通义千问聊天
// @Tags Qianwen
// @Summary 与通义千问聊天
// @Description 发送消息给通义千问并获取回复
// @Accept application/json
// @Produce application/json
// @Param data body request.ChatRequest true "聊天消息"
// @Success 200 {object} response.Response{data=response.ChatResponse,msg=string} "聊天成功"
// @Router /qianwen/chat [post]
func (q *QianwenApi) Chat(c *gin.Context) {
	var req request.ChatRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := qianwenService.Chat(req)
	if err != nil {
		global.GVA_LOG.Error("调用通义千问失败!", zap.Error(err))
		response.FailWithMessage("调用通义千问失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "请求成功", c)
}


