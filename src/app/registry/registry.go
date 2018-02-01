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
	Monitor []string
	SysctlNames map[string]string
	WsChan RegChan
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

	Registry.WsChan = NewRegChan()

	Registry.initServices(rr.Services)

	Registry.Session = session.Init()

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

func (r *registry) initServices(rr map[string]string) {

	r.SysctlNames = rr

	for n, _ := range r.SysctlNames {

		r.Monitor = append(r.Monitor, n)

	}

}