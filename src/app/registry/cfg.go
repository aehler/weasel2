package registry

import (
	"app/db"
)

type Cfg struct{
	Db map[string]db.Dbcreds
	Redis *Rediscreds
	Path *Path
	Services map[string]string
	Metrics []string
}