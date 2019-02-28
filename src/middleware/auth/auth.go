package auth

import (
	"app"
	"app/session"
	"app/registry"
	"lib/auth"
	"encoding/json"
	"fmt"
	"app/crypto"
)

type Auth struct {
	User *auth.User
	SSID string
}

func GetAuthUser(c *app.Context) {

	var sd string

	u := auth.Auth{
		User: &auth.User{},
	}

	//Get cookie
	err := session.Get(c.Request, &sd, &session.Config{Keys : registry.Registry.SessionKeys})

	if err == nil { //cookie was set

		//Get userData from registry
		v, err := registry.Registry.Session.Get(sd)

		if err != nil { //user is not authorized, using guest params

			u.User.UserID = 0
			u.User.UserLastName = "Unautharized"

			if err := registry.Registry.Session.Add(sd, u); err != nil {

				c.RenderJSONError(err)

				return
			}

		} else { //User is authorized, unmarshal data

			if err := json.Unmarshal(v, &u); err != nil {

				fmt.Println(err)

				return
			}

			u.User.SessionID = sd

			u.SSID = sd

		}

	} else { //cookie was not set, setting new cookie

		u.User.UserID = 0
		u.User.UserLastName = "Not logged in"

		ssid := crypto.GenSessionId(u.User.UserID, u.User.UserLastName)

		if err := session.Set(c.ResponseWriter, ssid, &session.Config{Keys: registry.Registry.SessionKeys}); err != nil {

			fmt.Println("couldn't set cookie")

			c.RenderJSONError(err)

			return

		}

		u.User.SessionID = ssid

		u.SSID = ssid

		if err := registry.Registry.Session.Add(ssid, u); err != nil {

			c.RenderJSONError(err)

			return
		}

		app.Redirect("/"+c.Request.URL.Path, c, 302)

		return

	}

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