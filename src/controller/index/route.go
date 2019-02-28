package index

import (
	"app"
	"middleware/auth"
	"middleware/user_settings"
)

func Route(ap *app.App) {

	ap.GetPost("/", auth.GetAuthUser, user_settings.GetUserSettings, Index)

	ap.Get("/dashboard/", auth.GetAuthUser, auth.Check, Dashboard)

}