package db

import (
	"github.com/jmoiron/sqlx"
	"log"
	"gopkg.in/mgo.v2"
)

type SQLConn interface {
	ConnectString() string
	Connect()
	SQLX() *sqlx.DB
	MGO() *mgo.Database
}

type Dbcreds struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func New(rr map[string]Dbcreds) SQLConn {

	for k, resource := range rr {

		switch k {

		case "postgresql":
			log.Println("Using Postgres")

			r := &PostgreSQL{
				Address : resource.Address,
				Username : resource.Username,
				Password : resource.Password,
				Database : resource.Database,
			}

			r.Connect()

			return r

		case "mysql":
			log.Println("Using MySQL")

			r := &MySQL{
				Address : resource.Address,
				Username : resource.Username,
				Password : resource.Password,
				Database : resource.Database,
			}

			r.Connect()

			return r

		case "mongodb":
			log.Println("Using MongoDB")

			r := &MongoDB{
				Address : resource.Address,
				Username : resource.Username,
				Password : resource.Password,
				Database : resource.Database,
			}

			r.Connect()

			return r

		}

	}

	return nil

}
