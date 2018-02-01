package guest

import (
	"fmt"
	"app/crypto"
	"app"
	"app/session"
	"app/registry"
	"lib/auth"
)

func GuestSettings(c *app.Context) {

	gc := auth.Auth{}

	if err := session.Get(c.Request, &gc, &session.Config{Keys : registry.Registry.SessionKeys, Name:"guest"}); err == nil {

		c.Set("ssid", gc.SSID)
		c.Set("lang", gc.Lang)

	}

	ssid := crypto.GenSessionId(0, "guest")

	gc.SSID = ssid
	gc.Lang = "en"

	if err := session.Set(c.ResponseWriter, gc, &session.Config{Keys : registry.Registry.SessionKeys, Name:"guest"}); err != nil {

		fmt.Println("couldn't set cookie")

		c.RenderError(err.Error())

		c.Stop()

		return

	}

	c.Set("ssid", gc.SSID)
	c.Set("lang", gc.Lang)

}

func ResetLanguage(c *app.Context) {

	gc := auth.Auth{}

	lang := c.Params.ByName("lang")

	if err := session.Get(c.Request, &gc, &session.Config{Keys : registry.Registry.SessionKeys}); err == nil {

		gc.Lang = lang

		if err := session.Set(c.ResponseWriter, gc, &session.Config{Keys : registry.Registry.SessionKeys}); err != nil {

			fmt.Println("couldn't set cookie")

			c.RenderError(err.Error())

			c.Stop()

			return

		}

		c.Set("lang", gc.Lang)

	} else {

		fmt.Println(err)

	}

}