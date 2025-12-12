package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/lottery"

// CheckLotteryReq 检查彩票中奖请求
type CheckLotteryReq struct {
	LotteryType  lottery.LotteryType `json:"lotteryType" binding:"required"`  // 彩票类型: 双色球 或 大乐透
	WinningSSQ   *lottery.SSQNumbers `json:"winningSSQ"`   // 双色球中奖号码
	WinningDLT   *lottery.DLTNumbers `json:"winningDLT"`   // 大乐透中奖号码
	PurchasedSSQ []lottery.SSQNumbers `json:"purchasedSSQ"` // 购买的双色球号码（批量）
	PurchasedDLT []lottery.DLTNumbers `json:"purchasedDLT"` // 购买的大乐透号码（批量）
}

