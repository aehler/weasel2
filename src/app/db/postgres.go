package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
)

type PostgreSQL struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	conn *sqlx.DB
}


func (p *PostgreSQL) ConnectString() string {

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&connect_timeout=1", p.Username, p.Password, p.Address, p.Database)
}

func (p *PostgreSQL) Connect() {

	p.conn = sqlx.MustOpen("postgres", p.ConnectString())

}

func (p *PostgreSQL) SQLX() *sqlx.DB {

	return p.conn

}

func (p *PostgreSQL) MGO() *mgo.Database {

	return nil

}