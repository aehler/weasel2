package personal

import (
	"app"
	"middleware/auth"
	"middleware/user_settings"
)

func Route(ap *app.App) {

	ap.Get("/my-blueprints/", auth.GetAuthUser, user_settings.GetUserSettings, myBPO)
	ap.Post("/pinned-report/", auth.GetAuthUser, user_settings.GetUserSettings, pinnedBPOReport)
	ap.Post("/toggle-my-blueprint/", auth.GetAuthUser, user_settings.GetUserSettings, togglePinBPO)

}