package logs

import "app/registry"

var Logs Logger

func Init() {

	if registry.Registry.Connect.SQLX() != nil {

		Logs = &sqlx_l{}

	} else if registry.Registry.Connect.MGO() != nil {

		Logs = &mgo_l{}

	} else {

		panic("Dunno which logs to use, none of the SQLX or MGO DB drivers inited!")

	}

}