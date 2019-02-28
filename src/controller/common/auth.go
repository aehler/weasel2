package common

import (
	"app"
	au "lib/auth"
	"app/registry"
	"middleware/auth"

)

func authCallback (c *app.Context){

	user := c.Get("user").(*au.User)

	c.Request.ParseForm()

	if c.Request.FormValue("login") == "" || c.Request.FormValue("password") == "" {
		c.RenderJSON(map[string]interface{}{
			"error": "Логин или пароль пусты",
			"status" : false,
		})
	}

	user, err := au.AuthUser(c.Request.FormValue("login"), c.Request.FormValue("password"))
	if err != nil {
		c.RenderJSONError(err)
		return
	}

	registry.Registry.Session.Replace(user.SessionID, auth.Auth{User: user, SSID: user.SessionID})

	c.Redirect("/dashboard/")

	return

}