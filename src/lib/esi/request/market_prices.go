package request

type MarketPrices struct{
	Path string
}

func (m MarketPrices) IsValid() bool {
	return true
}

func (m MarketPrices) Url() string {

	return m.Path

}

func (m MarketPrices) RequiresAuth() bool {
	return false
}

func (m MarketPrices) GetToken() string {
	return ""
}