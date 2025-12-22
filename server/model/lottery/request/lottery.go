package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/lottery"

// CheckLotteryReq 检查彩票中奖请求（支持批量多期）
type CheckLotteryReq struct {
	LotteryType  lottery.LotteryType   `json:"lotteryType" binding:"required"` // 彩票类型: 双色球 或 大乐透
	WinningSSQ   []lottery.SSQNumbers  `json:"winningSSQ"`                     // 双色球中奖号码（支持多期）
	WinningDLT   []lottery.DLTNumbers  `json:"winningDLT"`                     // 大乐透中奖号码（支持多期）
}

