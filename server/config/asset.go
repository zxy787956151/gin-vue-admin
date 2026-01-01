package config

type Asset struct {
	Items []AssetItem `mapstructure:"items" json:"items" yaml:"items"`
}

type AssetItem struct {
	Name  string  `mapstructure:"name" json:"name" yaml:"name"`
	Value float64 `mapstructure:"value" json:"value" yaml:"value"`
}
