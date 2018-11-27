package index

import (
	"app"
	"middleware/auth"
	"middleware/user_settings"
)

func Route(ap *app.App) {

	ap.GetPost("/", auth.GetAuthUser, user_settings.GetUserSettings, Index)

	ap.Get("/dashboard/", auth.GetAuthUser, auth.Check, Dashboard)

	ap.GetPost("/login/", func(c *app.Context){

		if c.IsPost(){

			c.RenderJSON(map[string]interface{}{
				"error": "Not yet implemented",
				"status" : false,
			})

			return
		}

		c.RenderHTML("/login.html", map[string]interface {} {

		})
	})

}