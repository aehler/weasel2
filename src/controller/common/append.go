package common

import (
	"app"
	"lib/common"
	"strconv"
	"lib/auth"
)

func appendUserSettings(c *app.Context) {

	user := c.Get("user").(*auth.User)

	if c.IsPost() {

		c.Request.ParseForm()

		data := common.UserSettings{}

		for key, _ := range c.Request.Form {

			switch key {

			case "regionID":

				val, err := strconv.ParseUint(c.Request.FormValue(key), 10, 64)
				if err != nil {

					c.RenderJSON(map[string]interface{}{
						"Error" : err.Error(),
						"Result": false,
					})

				}

				data.RegionID = uint(val)

			case "region":

				data.Region = c.Request.FormValue(key)

			case "settings-pi":

				val, err := strconv.ParseUint(c.Request.FormValue(key), 10, 64)
				if err != nil {

					c.RenderJSON(map[string]interface{}{
						"Error" : err.Error(),
						"Result": false,
					})

				}

				data.MarketShare = uint(val)

			}
		}

		if err := data.Save(user); err != nil {
			c.RenderJSON(map[string]interface{}{
				"Error" : err.Error(),
				"Result": false,
			})
		} else {
			c.RenderJSON(map[string]interface{}{
				"Error" : nil,
				"Result": true,
			})
		}

	} else {

		c.RenderJSON(map[string]interface{}{
			"Error" : "Bad request",
			"Result": false,
		})

	}

}