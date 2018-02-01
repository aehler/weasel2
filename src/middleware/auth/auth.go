package auth

import (
	"app"
	"app/session"
	"app/registry"
	"lib/auth"
	"encoding/json"
	"fmt"
)

func GetAuthUser(c *app.Context) {

	var sd string

	u := auth.Auth{}

	if err := session.Get(c.Request, &sd, &session.Config{Keys : registry.Registry.SessionKeys, Name:"auth"}); err != nil {

		fmt.Println(err)

		return
	}

	v, err := registry.Registry.Session.Get(sd)

	if err != nil {

		fmt.Println(err)

		return
	}

	if err := json.Unmarshal(v, &u); err != nil {

		fmt.Println(err)

		return
	}

	u.SSID = sd

	c.Set("user", u.User)

}

func Check(c *app.Context) {

	var sd string

	if err := session.Get(c.Request, &sd, &session.Config{Keys : registry.Registry.SessionKeys, Name:"auth"}); err != nil {

		fmt.Println(err)

		app.Redirect("/login/", c, 302)

		return
	}

	v, err := registry.Registry.Session.Get(sd)

	if err != nil {

		fmt.Println(err)

		app.Redirect("/login/", c, 302)

		return
	}

	u := auth.Auth{}

	if err := json.Unmarshal(v, &u); err != nil {

		fmt.Println(err)

		app.Redirect("/login/", c, 302)

		return
	}

	u.SSID = sd

	c.Set("user", u)
}