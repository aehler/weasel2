package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	"time"
)

type MongoDB struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	conn *mgo.Database
}

func (p *MongoDB) ConnectString() string {

	if p.Username == "" || p.Password == "" {
		return fmt.Sprintf("%s", p.Address)
	}

	return fmt.Sprintf("%s:%s@%s", p.Username, p.Password, p.Address)
}

func (p *MongoDB) Connect() {

	var (
		err error
		session *mgo.Session
	)

	log.Println("Attempting connection to MongoDB...")

	info := &mgo.DialInfo{
		Username: p.Username,
		Password: p.Password,
		Database: p.Database,
		Timeout: 1 * time.Second,
		Addrs: []string{p.Address},
	}

	if session, err = mgo.DialWithInfo(info); err != nil {

		log.Fatal("Connection to MongoDB failed!", err.Error())

	}

	session.SetMode(mgo.Monotonic, true)

	p.conn = session.DB(p.Database)

	log.Println("Connected to MongoDB.")

}

func (p *MongoDB) SQLX() *sqlx.DB {

	return nil

}

func (p *MongoDB) MGO() *mgo.Database {

	return p.conn

}