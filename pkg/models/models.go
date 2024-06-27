package models

type CryptoHistory struct {
	Rate      float64 `json:"rate" gorm:"column:rate"`
	Volume    float64 `json:"volume" gorm:"column:volume"`
	OrderType string  `json:"order_type" gorm:"column:order_type"`
	Timestamp string  `json:"timestamp" gorm:"column:timestamp;type:timestamp"`
}

type RateStats struct {
	MaxRate   float64 `json:"max_Rate"`
	MinRate   float64 `json:"min_rate"`
	SumVolume float64 `json:"sum_volume"`
}
