package registry

import (
	"app/session"
	"app/db"
	"github.com/KunBetter/MemKV"
)

type registry struct {
	Config Cfg
	Connect db.SQLConn
	Redis *Redis
	Session *session.SessionStorage
	SessionKeys []*[32]byte
	ReferenceConf map[string]*refConf
	storage map[string]*storage
	KVS *MemKV.MemKV
}

type Rediscreds struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DB     int    `yaml:"db"`
	MQName string `yaml:"mqname"`
	Password string
}

var Registry registry

func Init(rr Cfg, config string) {

	Registry.Config = rr

	Registry.KVS = MemKV.DB()

	Registry.Connect = db.New(rr.Db)

	Registry.Session = session.Init(rr.Memcached, rr.SessionTimeout)

	Registry.SessionKeys = append(Registry.SessionKeys, &[32]byte{
			'm',
		} )

	readRefConf(config)

	readStorageConf(config)

	Registry.newRedisClient(rr.Redis)

}

func (r *registry) Storage(key string) *storage {

	return Registry.storage[key]

}