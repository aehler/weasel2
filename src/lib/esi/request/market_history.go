package request

import "fmt"

type MarketHistory struct{
	Path string
	RegionID uint
	ItemID uint
}

func (m MarketHistory) IsValid() bool {
	return m.RegionID != 0 && m.ItemID != 0
}

func (m MarketHistory) Url() string {

	return fmt.Sprintf("markets/%d/history/?type_id=%d", m.RegionID, m.ItemID)

}

func (m MarketHistory) RequiresAuth() bool {
	return false
}

func (m MarketHistory) GetToken() string {
	return ""
}