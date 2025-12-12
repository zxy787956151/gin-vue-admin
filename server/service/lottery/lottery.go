package lottery

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery/response"
)

type LotteryService struct{}

// CheckLottery 检查彩票中奖情况
func (l *LotteryService) CheckLottery(req request.CheckLotteryReq) (response.CheckLotteryResp, error) {
	resp := response.CheckLotteryResp{
		LotteryType: req.LotteryType,
		Results:     []response.LotteryResult{},
	}

	switch req.LotteryType {
	case lottery.SSQ:
		// 双色球
		if req.WinningSSQ == nil {
			return resp, errors.New("双色球类型必须提供中奖号码")
		}
		if len(req.PurchasedSSQ) == 0 {
			return resp, errors.New("请提供购买的双色球号码")
		}
		// 验证中奖号码
		if err := validateSSQNumbers(*req.WinningSSQ); err != nil {
			return resp, err
		}
		// 批量检查
		for _, purchased := range req.PurchasedSSQ {
			if err := validateSSQNumbers(purchased); err != nil {
				return resp, err
			}
			result := l.checkSSQ(*req.WinningSSQ, purchased)
			resp.Results = append(resp.Results, result)
		}

	case lottery.DLT:
		// 大乐透
		if req.WinningDLT == nil {
			return resp, errors.New("大乐透类型必须提供中奖号码")
		}
		if len(req.PurchasedDLT) == 0 {
			return resp, errors.New("请提供购买的大乐透号码")
		}
		// 验证中奖号码
		if err := validateDLTNumbers(*req.WinningDLT); err != nil {
			return resp, err
		}
		// 批量检查
		for _, purchased := range req.PurchasedDLT {
			if err := validateDLTNumbers(purchased); err != nil {
				return resp, err
			}
			result := l.checkDLT(*req.WinningDLT, purchased)
			resp.Results = append(resp.Results, result)
		}

	default:
		return resp, errors.New("不支持的彩票类型，仅支持：双色球、大乐透")
	}

	return resp, nil
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

// checkSSQ 检查双色球中奖情况
func (l *LotteryService) checkSSQ(winning, purchased lottery.SSQNumbers) response.LotteryResult {
	// 计算红球匹配数
	redMatch := 0
	for _, pb := range purchased.RedBalls {
		for _, wb := range winning.RedBalls {
			if pb == wb {
				redMatch++
				break
			}
		}
	}

	// 蓝球是否匹配
	blueMatch := purchased.BlueBall == winning.BlueBall

	// 判断奖项
	prizeLevel := 0
	prize := "未中奖"
	matchDetail := fmt.Sprintf("红球匹配%d个", redMatch)
	if blueMatch {
		matchDetail += "，蓝球匹配"
	} else {
		matchDetail += "，蓝球未匹配"
	}

	if redMatch == 6 && blueMatch {
		prizeLevel = 1
		prize = "一等奖"
	} else if redMatch == 6 {
		prizeLevel = 2
		prize = "二等奖"
	} else if redMatch == 5 && blueMatch {
		prizeLevel = 3
		prize = "三等奖"
	} else if redMatch == 5 || (redMatch == 4 && blueMatch) {
		prizeLevel = 4
		prize = "四等奖"
	} else if redMatch == 4 || (redMatch == 3 && blueMatch) {
		prizeLevel = 5
		prize = "五等奖"
	} else if blueMatch {
		prizeLevel = 6
		prize = "六等奖"
	}

	return response.LotteryResult{
		Numbers:     purchased,
		IsWinning:   prizeLevel > 0,
		Prize:       prize,
		PrizeLevel:  prizeLevel,
		MatchDetail: matchDetail,
	}
}

// checkDLT 检查大乐透中奖情况
func (l *LotteryService) checkDLT(winning, purchased lottery.DLTNumbers) response.LotteryResult {
	// 计算前区匹配数
	frontMatch := 0
	for _, pf := range purchased.FrontZone {
		for _, wf := range winning.FrontZone {
			if pf == wf {
				frontMatch++
				break
			}
		}
	}

	// 计算后区匹配数
	backMatch := 0
	for _, pb := range purchased.BackZone {
		for _, wb := range winning.BackZone {
			if pb == wb {
				backMatch++
				break
			}
		}
	}

	// 判断奖项
	prizeLevel := 0
	prize := "未中奖"
	matchDetail := fmt.Sprintf("前区匹配%d个，后区匹配%d个", frontMatch, backMatch)

	if frontMatch == 5 && backMatch == 2 {
		prizeLevel = 1
		prize = "一等奖"
	} else if frontMatch == 5 && backMatch == 1 {
		prizeLevel = 2
		prize = "二等奖"
	} else if frontMatch == 5 && backMatch == 0 {
		prizeLevel = 3
		prize = "三等奖"
	} else if frontMatch == 4 && backMatch == 2 {
		prizeLevel = 4
		prize = "四等奖"
	} else if frontMatch == 4 && backMatch == 1 {
		prizeLevel = 5
		prize = "五等奖"
	} else if frontMatch == 3 && backMatch == 2 {
		prizeLevel = 6
		prize = "六等奖"
	} else if frontMatch == 4 && backMatch == 0 {
		prizeLevel = 7
		prize = "七等奖"
	} else if (frontMatch == 3 && backMatch == 1) || (frontMatch == 2 && backMatch == 2) {
		prizeLevel = 8
		prize = "八等奖"
	} else if (frontMatch == 3 && backMatch == 0) || (frontMatch == 1 && backMatch == 2) || (frontMatch == 2 && backMatch == 1) || (frontMatch == 0 && backMatch == 2) {
		prizeLevel = 9
		prize = "九等奖"
	}

	return response.LotteryResult{
		Numbers:     purchased,
		IsWinning:   prizeLevel > 0,
		Prize:       prize,
		PrizeLevel:  prizeLevel,
		MatchDetail: matchDetail,
	}
}

