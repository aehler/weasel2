package index

import (
	"app"
	"app/crypto"
	"lib/items"
	"lib/common"
	"lib/auth"
	"lib/esi/response"
	"app/registry"
)

func item(c *app.Context) {

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	settings := c.Get("userSettings").(common.UserSettings)

	user := c.Get("user").(*auth.User)

	item, err := items.GetItem(itemID, settings.RegionID, user)

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : "Get item error: " + err.Error(),
		})

		c.Stop()

		return
	}

	c.RenderHTML("/item-general.html", map[string]interface {} {
		"item"    : item,
		"selected" : "index",
	})

}

func itemManufacturing(c *app.Context) {

	settings := c.Get("userSettings").(common.UserSettings)

	user := c.Get("user").(*auth.User)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	item, err := items.GetItem(itemID, settings.RegionID, user)

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : "Get item error: " + err.Error(),
		})

		c.Stop()

		return
	}

	charSkills := response.CharacterSkills{}

	registry.Registry.Session.Unmarshal(user.SessionID+"_skills", &charSkills)

	c.RenderHTML("/item-manufacturing.html", map[string]interface {} {
		"item"    : item,
		"selected" : "manufacturing",
		"charSkills": charSkills.Skills,
	})

}

func itemResearch(c *app.Context) {

	user := c.Get("user").(*auth.User)

	settings := c.Get("userSettings").(common.UserSettings)

	itemID, err := crypto.DecryptUrl(c.Param("itemId"))

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	item, err := items.GetItem(itemID, settings.RegionID, user)

	if err != nil {
		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	c.RenderHTML("/item-research.html", map[string]interface {} {
		"item"    : item,
		"selected" : "research",
	})

}