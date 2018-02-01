package index

import (
	"app"
)

func Route(ap *app.App) {

	ap.Get("/", Index)
	ap.Get("/cpu/", CPU)
	ap.Get("/pids/", PIDS)
	ap.Get("/log/:sname/", Logs)
	ap.Get("/message/:logId/", LogDetails)
	ap.Get("/restart/:service/", restartService)

}