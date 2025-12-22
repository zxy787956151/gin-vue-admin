package lottery

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery/response"
)

type LotteryService struct{}

// 写死的购买号码
var (
	// 双色球购买号码: 红球 01 05 09 15 19 30, 蓝球 05
	mySSQNumbers = lottery.SSQNumbers{
		RedBalls: []int{1, 5, 9, 15, 19, 30},
		BlueBall: 5,
	}

	// 大乐透购买号码: 前区 05 09 15 19 30, 后区 02 10
	myDLTNumbers = lottery.DLTNumbers{
		FrontZone: []int{5, 9, 15, 19, 30},
		BackZone:  []int{2, 10},
	}
)

// CheckLottery 检查彩票中奖情况（支持批量多期）
func (l *LotteryService) CheckLottery(req request.CheckLotteryReq) (response.CheckLotteryResp, error) {
	resp := response.CheckLotteryResp{
		LotteryType: req.LotteryType,
		Results:     []response.PeriodResult{},
	}

	switch req.LotteryType {
	case lottery.SSQ:
		// 双色球
		if len(req.WinningSSQ) == 0 {
			return resp, errors.New("双色球类型必须提供至少一期中奖号码")
		}
		
		// 批量检测每一期
		for i, winning := range req.WinningSSQ {
			// 验证中奖号码
			if err := validateSSQNumbers(winning); err != nil {
				return resp, fmt.Errorf("第%d期号码格式错误: %v", i+1, err)
			}
			
			// 检测本期
			result := l.checkSSQPeriod(winning, i+1)
			resp.Results = append(resp.Results, result)
			
			if result.IsWinning {
				resp.WinningCount++
			}
		}
		
		resp.TotalPeriods = len(req.WinningSSQ)

	case lottery.DLT:
		// 大乐透
		if len(req.WinningDLT) == 0 {
			return resp, errors.New("大乐透类型必须提供至少一期中奖号码")
		}
		
		// 批量检测每一期
		for i, winning := range req.WinningDLT {
			// 验证中奖号码
			if err := validateDLTNumbers(winning); err != nil {
				return resp, fmt.Errorf("第%d期号码格式错误: %v", i+1, err)
			}
			
			// 检测本期
			result := l.checkDLTPeriod(winning, i+1)
			resp.Results = append(resp.Results, result)
			
			if result.IsWinning {
				resp.WinningCount++
			}
		}
		
		resp.TotalPeriods = len(req.WinningDLT)

	default:
		return resp, errors.New("不支持的彩票类型，仅支持：双色球、大乐透")
	}

	return resp, nil
}

// checkSSQPeriod 检查单期双色球
func (l *LotteryService) checkSSQPeriod(winning lottery.SSQNumbers, period int) response.PeriodResult {
	result := response.PeriodResult{
		Period:           period,
		PurchasedNumbers: mySSQNumbers,
		WinningNumbers:   winning,
	}
	
	// 计算红球匹配数
	redMatch := 0
	for _, pb := range mySSQNumbers.RedBalls {
		for _, wb := range winning.RedBalls {
			if pb == wb {
				redMatch++
				break
			}
		}
	}
	
	// 蓝球是否匹配
	blueMatch := mySSQNumbers.BlueBall == winning.BlueBall
	
	// 判断奖项
	result.MatchDetail = fmt.Sprintf("红球匹配%d个", redMatch)
	if blueMatch {
		result.MatchDetail += "，蓝球匹配"
	} else {
		result.MatchDetail += "，蓝球未匹配"
	}
	
	if redMatch == 6 && blueMatch {
		result.PrizeLevel = 1
		result.Prize = "一等奖"
		result.IsWinning = true
	} else if redMatch == 6 {
		result.PrizeLevel = 2
		result.Prize = "二等奖"
		result.IsWinning = true
	} else if redMatch == 5 && blueMatch {
		result.PrizeLevel = 3
		result.Prize = "三等奖"
		result.IsWinning = true
	} else if redMatch == 5 || (redMatch == 4 && blueMatch) {
		result.PrizeLevel = 4
		result.Prize = "四等奖"
		result.IsWinning = true
	} else if redMatch == 4 || (redMatch == 3 && blueMatch) {
		result.PrizeLevel = 5
		result.Prize = "五等奖"
		result.IsWinning = true
	} else if blueMatch {
		result.PrizeLevel = 6
		result.Prize = "六等奖"
		result.IsWinning = true
	} else {
		result.PrizeLevel = 0
		result.Prize = "未中奖"
		result.IsWinning = false
	}
	
	return result
}

// checkDLTPeriod 检查单期大乐透
func (l *LotteryService) checkDLTPeriod(winning lottery.DLTNumbers, period int) response.PeriodResult {
	result := response.PeriodResult{
		Period:           period,
		PurchasedNumbers: myDLTNumbers,
		WinningNumbers:   winning,
	}
	
	// 计算前区匹配数
	frontMatch := 0
	for _, pf := range myDLTNumbers.FrontZone {
		for _, wf := range winning.FrontZone {
			if pf == wf {
				frontMatch++
				break
			}
		}
	}
	
	// 计算后区匹配数
	backMatch := 0
	for _, pb := range myDLTNumbers.BackZone {
		for _, wb := range winning.BackZone {
			if pb == wb {
				backMatch++
				break
			}
		}
	}
	
	// 判断奖项
	result.MatchDetail = fmt.Sprintf("前区匹配%d个，后区匹配%d个", frontMatch, backMatch)
	
	if frontMatch == 5 && backMatch == 2 {
		result.PrizeLevel = 1
		result.Prize = "一等奖"
		result.IsWinning = true
	} else if frontMatch == 5 && backMatch == 1 {
		result.PrizeLevel = 2
		result.Prize = "二等奖"
		result.IsWinning = true
	} else if frontMatch == 5 && backMatch == 0 {
		result.PrizeLevel = 3
		result.Prize = "三等奖"
		result.IsWinning = true
	} else if frontMatch == 4 && backMatch == 2 {
		result.PrizeLevel = 4
		result.Prize = "四等奖"
		result.IsWinning = true
	} else if frontMatch == 4 && backMatch == 1 {
		result.PrizeLevel = 5
		result.Prize = "五等奖"
		result.IsWinning = true
	} else if frontMatch == 3 && backMatch == 2 {
		result.PrizeLevel = 6
		result.Prize = "六等奖"
		result.IsWinning = true
	} else if frontMatch == 4 && backMatch == 0 {
		result.PrizeLevel = 7
		result.Prize = "七等奖"
		result.IsWinning = true
	} else if (frontMatch == 3 && backMatch == 1) || (frontMatch == 2 && backMatch == 2) {
		result.PrizeLevel = 8
		result.Prize = "八等奖"
		result.IsWinning = true
	} else if (frontMatch == 3 && backMatch == 0) || (frontMatch == 1 && backMatch == 2) || (frontMatch == 2 && backMatch == 1) || (frontMatch == 0 && backMatch == 2) {
		result.PrizeLevel = 9
		result.Prize = "九等奖"
		result.IsWinning = true
	} else {
		result.PrizeLevel = 0
		result.Prize = "未中奖"
		result.IsWinning = false
	}
	
	return result
}

// validateSSQNumbers 验证双色球号码
func validateSSQNumbers(numbers lottery.SSQNumbers) error {
	if len(numbers.RedBalls) != 6 {
		return errors.New("双色球红球必须是6个")
	}
	for _, ball := range numbers.RedBalls {
		if ball < 1 || ball > 33 {
			return fmt.Errorf("双色球红球号码必须在1-33之间，当前值: %d", ball)
		}
	}
	if numbers.BlueBall < 1 || numbers.BlueBall > 16 {
		return fmt.Errorf("双色球蓝球号码必须在1-16之间，当前值: %d", numbers.BlueBall)
	}
	// 检查红球是否有重复
	seen := make(map[int]bool)
	for _, ball := range numbers.RedBalls {
		if seen[ball] {
			return fmt.Errorf("双色球红球号码不能重复: %d", ball)
		}
		seen[ball] = true
	}
	return nil
}

// validateDLTNumbers 验证大乐透号码
func validateDLTNumbers(numbers lottery.DLTNumbers) error {
	if len(numbers.FrontZone) != 5 {
		return errors.New("大乐透前区必须是5个号码")
	}
	if len(numbers.BackZone) != 2 {
		return errors.New("大乐透后区必须是2个号码")
	}
	for _, num := range numbers.FrontZone {
		if num < 1 || num > 35 {
			return fmt.Errorf("大乐透前区号码必须在1-35之间，当前值: %d", num)
		}
	}
	for _, num := range numbers.BackZone {
		if num < 1 || num > 12 {
			return fmt.Errorf("大乐透后区号码必须在1-12之间，当前值: %d", num)
		}
	}
	// 检查前区是否有重复
	seen := make(map[int]bool)
	for _, num := range numbers.FrontZone {
		if seen[num] {
			return fmt.Errorf("大乐透前区号码不能重复: %d", num)
		}
		seen[num] = true
	}
	// 检查后区是否有重复
	seen = make(map[int]bool)
	for _, num := range numbers.BackZone {
		if seen[num] {
			return fmt.Errorf("大乐透后区号码不能重复: %d", num)
		}
		seen[num] = true
	}
	return nil
}

