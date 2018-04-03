package registry

import (
	"app/db"
)

type Cfg struct{
	Db map[string]db.Dbcreds
	Redis *Rediscreds
	Path *Path
	SessionTimeout int32 `yaml:"session-timeout"`
	Memcached string
}