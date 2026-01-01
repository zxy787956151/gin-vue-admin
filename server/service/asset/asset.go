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
	case "asset3":
		configItems = global.GVA_CONFIG.Asset3.Items
	case "asset4":
		configItems = global.GVA_CONFIG.Asset4.Items
	case "asset5":
		configItems = global.GVA_CONFIG.Asset5.Items
	case "asset6":
		configItems = global.GVA_CONFIG.Asset6.Items
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

// calculateChildSum 计算子配置的总和
func (s *AssetService) calculateChildSum(configType string) float64 {
	var configItems []config.AssetItem
	switch configType {
	case "asset3":
		configItems = global.GVA_CONFIG.Asset3.Items
	case "asset4":
		configItems = global.GVA_CONFIG.Asset4.Items
	case "asset5":
		configItems = global.GVA_CONFIG.Asset5.Items
	case "asset6":
		configItems = global.GVA_CONFIG.Asset6.Items
	default:
		return 0
	}
	
	var sum float64
	for _, item := range configItems {
		sum += item.Value
	}
	return sum
}

// GetAssetDistribution 获取个人资产分布数据（使用asset配置）
// 对于有子配置的父级项，自动计算子配置的总和作为父级项的值
func (s *AssetService) GetAssetDistribution() (data AssetDistributionResponse, err error) {
	// 获取父级配置项（asset配置）
	parentItems := global.GVA_CONFIG.Asset.Items
	
	// 配置项名称到子配置类型的映射
	itemNameToChildConfigType := map[string]string{
		"活钱管理": "asset3",
		"稳健理财": "asset4",
		"长期投资": "asset5",
		"保险保障": "asset6",
	}
	
	items := make([]AssetDistributionItem, len(parentItems))
	var total float64
	
	for i, parentItem := range parentItems {
		// 检查是否有对应的子配置
		if childConfigType, hasChild := itemNameToChildConfigType[parentItem.Name]; hasChild {
			// 有子配置：计算子配置的总和作为父级项的值
			items[i] = AssetDistributionItem{
				Name:  parentItem.Name,
				Value: s.calculateChildSum(childConfigType),
			}
		} else {
			// 没有子配置：使用配置中的值（如 asset2 的情况）
			items[i] = AssetDistributionItem{
				Name:  parentItem.Name,
				Value: parentItem.Value,
			}
		}
		total += items[i].Value
	}
	
	data = AssetDistributionResponse{
		Total:     total,
		Items:     items,
		Timestamp: time.Now().Unix(),
	}
	
	return data, nil
}

// GetAssetDistribution2 获取个人资产分布数据2（使用asset2配置）
func (s *AssetService) GetAssetDistribution2() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset2")
}

// GetAssetDistribution3 获取个人资产分布数据3（使用asset3配置）
func (s *AssetService) GetAssetDistribution3() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset3")
}

// GetAssetDistribution4 获取个人资产分布数据4（使用asset4配置）
func (s *AssetService) GetAssetDistribution4() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset4")
}

// GetAssetDistribution5 获取个人资产分布数据5（使用asset5配置）
func (s *AssetService) GetAssetDistribution5() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset5")
}

// GetAssetDistribution6 获取个人资产分布数据6（使用asset6配置）
func (s *AssetService) GetAssetDistribution6() (data AssetDistributionResponse, err error) {
	return s.getAssetDistributionByType("asset6")
}

// GetAssetDistributionByItemName 根据配置项名称获取对应的详情数据
// itemName: "活钱管理" -> asset3, "稳健理财" -> asset4, "长期投资" -> asset5, "保险保障" -> asset6
func (s *AssetService) GetAssetDistributionByItemName(itemName string) (data AssetDistributionResponse, err error) {
	// 配置项名称到配置类型的映射
	itemNameToConfigType := map[string]string{
		"活钱管理": "asset3",
		"稳健理财": "asset4",
		"长期投资": "asset5",
		"保险保障": "asset6",
	}
	
	configType, ok := itemNameToConfigType[itemName]
	if !ok {
		// 如果找不到对应的配置，返回asset3作为默认值
		configType = "asset3"
	}
	
	return s.getAssetDistributionByType(configType)
}
