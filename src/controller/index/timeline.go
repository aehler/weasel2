package index

import (
	"app"
	"app/crypto"
	"lib/common"
	"lib/items"
	"strconv"
	"lib/auth"
)

func timeline(c *app.Context) {

	settings := c.Get("userSettings").(common.UserSettings)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	user := c.Get("user").(*auth.User)

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	itemData, err := items.GetItem(itemID, settings.RegionID, user)
	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	c.RenderHTML("/timeline.html", map[string]interface {} {
		"item"      : itemData,
		"Region"    : settings.Region,
		"selected"  : "timeline",
		"prodBatch" : itemData.Batch,
	})

}

func timelineData(c *app.Context) {

	settings := c.Get("userSettings").(common.UserSettings)

	user := c.Get("user").(*auth.User)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))
	if err != nil {
		c.RenderJSON(map[string]interface {} {
			"Error" : err.Error(),
			"Result" : nil,
		})

		c.Stop()

		return
	}

	var batch uint = 0

	if c.GetUrlParam("batch") != "" {
		b, err := strconv.ParseUint(c.GetUrlParam("batch"), 10, 64)
		if err != nil {
			batch = 1
		} else {
			batch = uint(b)
		}
	}


	ptl, err := items.ProductionTimeline(itemID, settings.RegionID, batch, user)
	if err != nil {
		c.RenderJSON(map[string]interface {} {
			"Error" : err.Error(),
			"Result" : nil,
		})

		c.Stop()

		return
	}

	c.RenderJSON(map[string]interface {} {
		"Error" : nil,
		"Result" : ptl,
	})

}