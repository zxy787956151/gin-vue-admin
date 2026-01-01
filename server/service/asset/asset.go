package asset

import (
	"time"
)

type AssetDistributionItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type AssetDistributionResponse struct {
	Total     float64                 `json:"total"`
	Items     []AssetDistributionItem `json:"items"`
	Timestamp int64                   `json:"timestamp"`
}

type AssetService struct{}

var AssetServiceApp = new(AssetService)

// GetAssetDistribution 获取个人资产分布数据（Mock数据）
func (s *AssetService) GetAssetDistribution() (data AssetDistributionResponse, err error) {
	// Mock数据：模拟不同类型的资产（固定数据，可以手动修改）
	items := []AssetDistributionItem{
		{
			Name:  "股票",
			Value: 250000.00, // 25万，可以修改这个值
		},
		{
			Name:  "基金",
			Value: 180000.00, // 18万，可以修改这个值
		},
		{
			Name:  "银行存款",
			Value: 120000.00, // 12万，可以修改这个值
		},
		{
			Name:  "理财产品",
			Value: 150000.00, // 15万，可以修改这个值
		},
		{
			Name:  "加密货币",
			Value: 80000.00, // 8万，可以修改这个值
		},
		{
			Name:  "房产",
			Value: 500000.00, // 50万，可以修改这个值
		},
	}

	// 计算总资产
	var total float64
	for _, item := range items {
		total += item.Value
	}

	data = AssetDistributionResponse{
		Total:     total, // 总资产：128万
		Items:     items,
		Timestamp: time.Now().Unix(),
	}

	return data, nil
}
