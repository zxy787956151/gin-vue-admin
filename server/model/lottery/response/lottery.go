package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/lottery"

// PeriodResult 单期检测结果
type PeriodResult struct {
	Period           int         `json:"period"`           // 期号（从1开始）
	PurchasedNumbers interface{} `json:"purchasedNumbers"` // 购买的号码
	WinningNumbers   interface{} `json:"winningNumbers"`   // 中奖号码
	IsWinning        bool        `json:"isWinning"`        // 是否中奖
	Prize            string      `json:"prize"`            // 奖项: 一等奖、二等奖、三等奖等
	PrizeLevel       int         `json:"prizeLevel"`       // 奖项等级: 1-一等奖, 2-二等奖, 0-未中奖
	MatchDetail      string      `json:"matchDetail"`      // 匹配详情
}

// CheckLotteryResp 检查彩票中奖响应（支持批量多期）
type CheckLotteryResp struct {
	LotteryType    lottery.LotteryType `json:"lotteryType"`    // 彩票类型
	TotalPeriods   int                 `json:"totalPeriods"`   // 总期数
	WinningCount   int                 `json:"winningCount"`   // 中奖期数
	Results        []PeriodResult      `json:"results"`        // 每期的检测结果
}

