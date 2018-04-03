package esi

import (
	"lib/esi/request"
	"lib/esi/response"
	"fmt"
)

func (c *Client) MarketPrices() (response.MarketPrices, error) {

	res := response.MarketPrices{}

	if err := c.rpc.Send(request.MarketPrices{Path:"markets/prices/"}, &res); err != nil {

		fmt.Println(err.Error())

		return res, err

	}

	return res, nil

}

func (c *Client) MarketHistory(regionID, itemID uint) (response.MarketHistory, error) {

	res := response.MarketHistory{}

	if err := c.rpc.Send(request.MarketHistory{ItemID: itemID, RegionID: regionID}, &res); err != nil {

		fmt.Println(err.Error())

		return res, err

	}

	return res, nil

}