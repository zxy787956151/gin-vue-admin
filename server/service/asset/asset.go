package asset

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
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

// getAssetDistributionByType 核心逻辑：根据配置类型获取资产分布数据
func (s *AssetService) getAssetDistributionByType(configType string) (data AssetDistributionResponse, err error) {
	// 根据配置类型选择不同的配置
	var configItems []config.AssetItem
	switch configType {
	case "asset2":
		configItems = global.GVA_CONFIG.Asset2.Items
	case "asset":
		fallthrough
	default:
		configItems = global.GVA_CONFIG.Asset.Items
	}
	
	items := make([]AssetDistributionItem, len(configItems))
	
	for i, configItem := range configItems {
		items[i] = AssetDistributionItem{
			Name:  configItem.Name,
			Value: configItem.Value,
		}
	}

	// 计算总资产
	var total float64
	for _, item := range items {
		total += item.Value
	}

	data = AssetDistributionResponse{
		Total:     total,
		Items:     items,
		Timestamp: time.Now().Unix(),
	}

	return data, nil
}

// GetAssetDistribution 获取个人资产分布数据（使用asset配置）
func (s *AssetService) GetAssetDistribution() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset")
}

// GetAssetDistribution2 获取个人资产分布数据2（使用asset2配置）
func (s *AssetService) GetAssetDistribution2() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset2")
}
