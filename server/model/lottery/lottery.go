package lottery

// LotteryType 彩票类型
type LotteryType string

const (
	// SSQ 双色球
	SSQ LotteryType = "双色球"
	// DLT 大乐透
	DLT LotteryType = "大乐透"
)

// SSQNumbers 双色球号码
type SSQNumbers struct {
	RedBalls  []int `json:"redBalls" binding:"required,len=6"`  // 红球6个: 1-33
	BlueBall  int   `json:"blueBall" binding:"required,min=1,max=16"` // 蓝球1个: 1-16
}

// DLTNumbers 大乐透号码
type DLTNumbers struct {
	FrontZone []int `json:"frontZone" binding:"required,len=5"` // 前区5个: 1-35
	BackZone  []int `json:"backZone" binding:"required,len=2"`  // 后区2个: 1-12
}

