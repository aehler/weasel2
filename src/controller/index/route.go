package index

import (
	"app"
	"middleware/auth"
	"middleware/user_settings"
)

func Route(ap *app.App) {

	ap.GetPost("/", auth.GetAuthUser, user_settings.GetUserSettings, Index)
	ap.Get("/page/:page/", auth.GetAuthUser, user_settings.GetUserSettings, Index)
	ap.Get("/item/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, item)
	ap.Get("/item-manufacturing/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, itemManufacturing)
	ap.Get("/item-research/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, itemResearch)
	ap.Get("/market/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, market)
	ap.Get("/market-history/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, marketHistory)
	ap.Get("/timeline/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, timeline)
	ap.Get("/timeline-data/:itemId/", auth.GetAuthUser, user_settings.GetUserSettings, timelineData)

}