package logs

import (
	"time"
	"github.com/akdcode/monitor/protocols"
	"gopkg.in/mgo.v2/bson"
)

type LogEntry struct {
	MGOObjectId bson.ObjectId `db:"-" bson:"_id,omitempty"`
	LogID string `db:"log_id" bson:"-"`
	Service string `db:"service" bson:"service"`
	SeverityLevel int `db:"severity_level" bson:"severity_level"`
	Info string `db:"info" bson:"info,omitempty"`
	Details string `db:"details" bson:"details,omitempty"`
	Occured time.Time `db:"occured" bson:"occured,omitempty"`
	Created time.Time `db:"created" bson:"created,omitempty"`
}

type Logger interface {
	GetServiceLogs(s string, limit, offset uint) ([]LogEntry, error)
	GetServiceLogEntry(id string) (LogEntry, error)
	Store(p *protocols.Message) error
}