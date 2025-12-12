package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/lottery"

// LotteryResult 单个彩票的中奖结果
type LotteryResult struct {
	Numbers     interface{} `json:"numbers"`     // 号码（SSQNumbers 或 DLTNumbers）
	IsWinning   bool        `json:"isWinning"`   // 是否中奖
	Prize       string      `json:"prize"`       // 奖项: 一等奖、二等奖、三等奖等
	PrizeLevel  int         `json:"prizeLevel"`  // 奖项等级: 1-一等奖, 2-二等奖, 0-未中奖
	MatchDetail string      `json:"matchDetail"` // 匹配详情
}

// CheckLotteryResp 检查彩票中奖响应
type CheckLotteryResp struct {
	LotteryType lottery.LotteryType `json:"lotteryType"` // 彩票类型
	Results     []LotteryResult     `json:"results"`     // 所有彩票的中奖结果
}

