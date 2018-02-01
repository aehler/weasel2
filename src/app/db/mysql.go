package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type MySQL struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	conn *sqlx.DB
}

func (p *MySQL) ConnectString() string {

	return fmt.Sprintf("%s:%s@%s/%s", p.Username, p.Password, p.Address, p.Database)
}

func (p *MySQL) Connect() {

	p.conn = sqlx.MustOpen("mysql", p.ConnectString())

}

func (p *MySQL) SQLX() *sqlx.DB {

	return p.conn

}

func (p *MySQL) MGO() *mgo.Database {

	return nil

}