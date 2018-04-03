package response

import (

)

type MarketHistory []MHItem

type MHItem struct{
	Date EsiDate `json:"date" db:"dt"`
	OrderCount uint `json:"order_count" db:"order_count"`
	Volume uint `json:"volume" db:"volume"`
	Highest float64 `json:"highest" db:"highest"`
	Average float64 `json:"average" db:"average"`
	Lowest float64 `json:"lowest" db:"lowest"`
}