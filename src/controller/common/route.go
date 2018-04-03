package common

import (
	"app"
	"lib/common"
	"middleware/auth"
	"app/crypto"
	au "lib/auth"
)

func Route(ap *app.App) {

	ap.Get("/list/regions/", func(c *app.Context){

		c.RenderJSON(common.ListRegions(c.GetUrlParam("term")))

	})

	ap.Post("/settings/append/", auth.GetAuthUser, appendUserSettings)

	ap.Get("/login/", auth.GetAuthUser, func(c *app.Context){

		user := c.Get("user").(*au.User)

		ccs := crypto.EncryptB64(user.SessionID, key)

		c.Redirect(`https://login.eveonline.com/oauth/authorize/?response_type=code&redirect_uri=https://127.0.0.1:8087/login-success/&client_id=bbacb081bc184538b1c8aa360036bd82&scope=characterSkillsRead esi-skills.read_skills.v1&state=`+ccs)

	})

	ap.Get("/login-success/", auth.GetAuthUser, oauthCallback)
}