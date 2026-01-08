package localai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/localai/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LocalAIApi struct{}

// Chat 本地AI聊天（支持RAG）
// @Tags LocalAI
// @Summary 本地AI聊天
// @Description 使用本地模型进行对话，支持RAG检索增强
// @Accept application/json
// @Produce application/json
// @Param data body request.ChatRequest true "聊天请求"
// @Success 200 {object} response.Response{data=response.ChatResponse,msg=string} "聊天成功"
// @Router /localai/chat [post]
func (l *LocalAIApi) Chat(c *gin.Context) {
	var req request.ChatRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := getService().Chat(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("本地AI聊天失败!", zap.Error(err))
		response.FailWithMessage("聊天失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "请求成功", c)
}

// IngestDocument 导入文档到知识库
// @Tags LocalAI
// @Summary 导入文档
// @Description 将文档内容添加到向量知识库
// @Accept application/json
// @Produce application/json
// @Param data body request.IngestRequest true "文档内容"
// @Success 200 {object} response.Response{data=response.IngestResponse,msg=string} "导入成功"
// @Router /localai/ingest [post]
func (l *LocalAIApi) IngestDocument(c *gin.Context) {
	var req request.IngestRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := getService().IngestDocument(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("导入文档失败!", zap.Error(err))
		response.FailWithMessage("导入失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "文档导入成功", c)
}

// SearchDocuments 搜索文档
// @Tags LocalAI
// @Summary 搜索文档
// @Description 在知识库中搜索相关文档
// @Accept application/json
// @Produce application/json
// @Param data body request.SearchRequest true "搜索请求"
// @Success 200 {object} response.Response{data=response.SearchResponse,msg=string} "搜索成功"
// @Router /localai/search [post]
func (l *LocalAIApi) SearchDocuments(c *gin.Context) {
	var req request.SearchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := getService().SearchDocuments(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("搜索文档失败!", zap.Error(err))
		response.FailWithMessage("搜索失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "搜索成功", c)
}

// SubmitFeedback 提交用户反馈
// @Tags LocalAI
// @Summary 提交反馈
// @Description 提交对AI回答的反馈，用于模型训练
// @Accept application/json
// @Produce application/json
// @Param data body request.FeedbackRequest true "反馈内容"
// @Success 200 {object} response.Response{msg=string} "提交成功"
// @Router /localai/feedback [post]
func (l *LocalAIApi) SubmitFeedback(c *gin.Context) {
	var req request.FeedbackRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = getService().SubmitFeedback(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("提交反馈失败!", zap.Error(err))
		response.FailWithMessage("提交失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("反馈已提交，感谢您的反馈！", c)
}

// StartTraining 开始训练模型
// @Tags LocalAI
// @Summary 开始训练
// @Description 基于收集的反馈数据训练模型
// @Accept application/json
// @Produce application/json
// @Param data body request.TrainRequest true "训练参数"
// @Success 200 {object} response.Response{msg=string} "训练已启动"
// @Router /localai/train/start [post]
func (l *LocalAIApi) StartTraining(c *gin.Context) {
	var req request.TrainRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 异步启动训练
	go func() {
		if err := getService().StartTraining(c.Request.Context(), req); err != nil {
			global.GVA_LOG.Error("训练失败!", zap.Error(err))
		}
	}()

	response.OkWithMessage("训练已启动，请稍后查询训练状态", c)
}

// GetTrainingStatus 获取训练状态
// @Tags LocalAI
// @Summary 获取训练状态
// @Description 查询当前的训练状态和进度
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.TrainStatusResponse,msg=string} "查询成功"
// @Router /localai/train/status [get]
func (l *LocalAIApi) GetTrainingStatus(c *gin.Context) {
	result := getService().GetTrainingStatus()
	response.OkWithDetailed(result, "查询成功", c)
}

// GetStats 获取系统统计信息
// @Tags LocalAI
// @Summary 获取统计信息
// @Description 获取知识库和训练数据的统计信息
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.StatsResponse,msg=string} "查询成功"
// @Router /localai/stats [get]
func (l *LocalAIApi) GetStats(c *gin.Context) {
	result := getService().GetStats()
	response.OkWithDetailed(result, "查询成功", c)
}

// HealthCheck 健康检查
// @Tags LocalAI
// @Summary 健康检查
// @Description 检查本地模型服务状态
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "服务正常"
// @Router /localai/health [get]
func (l *LocalAIApi) HealthCheck(c *gin.Context) {
	err := getService().CheckHealth(c.Request.Context())
	if err != nil {
		response.FailWithMessage("服务异常: "+err.Error(), c)
		return
	}
	response.OkWithMessage("服务正常", c)
}

// DeleteDocument 删除文档
// @Tags LocalAI
// @Summary 删除文档
// @Description 从知识库中删除指定文档
// @Accept application/json
// @Produce application/json
// @Param data body request.DeleteDocRequest true "文档ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /localai/document/delete [post]
func (l *LocalAIApi) DeleteDocument(c *gin.Context) {
	var req request.DeleteDocRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = getService().DeleteDocument(c.Request.Context(), req.ID)
	if err != nil {
		global.GVA_LOG.Error("删除文档失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("文档已删除", c)
}


