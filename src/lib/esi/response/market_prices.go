package response

type MarketPrices []AvgPrice

type AvgPrice struct{
	TypeID uint `json:"type_id"`
	AvgPrice float64 `json:"average_price"`
	AdjPrice float64 `json:"adjusted_price"`
}