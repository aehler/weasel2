package controller

import (
	"controller/index"
	"controller/common"
	"controller/personal"
	"app"
)

func Route(a *app.App) {
	index.Route(a)
	common.Route(a)
	personal.Route(a)
}
