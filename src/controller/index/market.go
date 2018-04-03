package index

import (
	"app"
	"app/crypto"
	m "lib/market"
	"lib/common"
	"lib/items"
	"fmt"
	"strconv"
)

func market(c *app.Context) {

	settings := c.Get("userSettings").(common.UserSettings)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	itemData, err := items.GetType(itemID)
	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	bpoData, err := items.GetTypeByBPO(itemID)
	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	itemData.BpoName = itemData.TypeName

	c.RenderHTML("/market.html", map[string]interface {} {
		"item"      : itemData,
		"itemData"  : bpoData,
		"Region"    : settings.Region,
		"selected" : "market",
	})

}

func marketHistory(c *app.Context) {

	settings := c.Get("userSettings").(common.UserSettings)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	md, err := m.GetMarketHistory(itemID, settings.RegionID)
	if err != nil {
		c.RenderJSON(map[string]interface {} {
			"Error" : err.Error(),
			"Result" : nil,
		})

		c.Stop()

		return
	}

	var hcData = [][]interface{}{}

	for _, mm := range md {

		ts := mm.Date.Unix()
		stamp, _ := strconv.Atoi(fmt.Sprint(ts))

		hcData = append(hcData, []interface{}{
			stamp*1000,
			mm.Average,
			mm.Highest,
			mm.Lowest,
			mm.Volume,
		})

	}

	c.RenderJSON(map[string]interface {} {
		"Error" : nil,
		"Result" : hcData,
	})

}