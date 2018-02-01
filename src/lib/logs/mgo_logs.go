package logs

import (
	"app/registry"
	"github.com/akdcode/monitor/protocols"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type mgo_l struct{}

func (_ *mgo_l) GetServiceLogs(s string, limit, offset uint) ([]LogEntry, error) {

	res := []LogEntry{}
	l := LogEntry{}

	i := registry.Registry.Connect.MGO().C("logs").Find(bson.M{
		"service" : s,
	}).Sort("-$natural").Limit(int(limit)).Skip(int(offset)).Iter()

	for i.Next(&l) {

		l.LogID = l.MGOObjectId.Hex()

		res = append(res, l)

	}

	return res, nil

}

func (_ *mgo_l) GetServiceLogEntry(id string) (LogEntry, error) {

	res := LogEntry{}

	err := registry.Registry.Connect.MGO().C("logs").FindId(bson.ObjectIdHex(id)).One(&res)

	if err != nil {
		return res, err
	}

	return res, nil

}

func (_ *mgo_l) Store(p *protocols.Message) error {

	details, _ := json.MarshalIndent(p.Data, "", "	")

	il := LogEntry{
		MGOObjectId: bson.NewObjectId(),
		Service: p.Owner,
		SeverityLevel: int(p.MessageType),
		Info: p.Info,
		Details: string(details),
		Occured: p.Time,
		Created: time.Now(),
	}

	if err := registry.Registry.Connect.MGO().C("logs").Insert(il); err != nil {

		return err

	}

	p.SetMonitorId(il.MGOObjectId.Hex())

	return nil

}