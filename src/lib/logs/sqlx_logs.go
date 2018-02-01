package logs

import (
	"app/registry"
	"github.com/akdcode/monitor/protocols"
	"encoding/json"
	"strconv"
)

type sqlx_l struct{}

func (_ *sqlx_l) GetServiceLogs(s string, limit, offset uint) ([]LogEntry, error) {

	res := []LogEntry{}

	if err := registry.Registry.Connect.SQLX().Select(&res, `select log_id, service, severity_level, info, '' as details, occured, created
	from service_logs.logs where service = $1 order by log_id desc limit $2 offset $3`,
		s,
		limit,
		offset,
		); err != nil {

			return res, err

	}

	return res, nil

}

func (_ *sqlx_l) GetServiceLogEntry(id string) (LogEntry, error) {

	res := LogEntry{}

	uid, err := strconv.Atoi(id)
	if err != nil {

		return res, err

	}

	if err := registry.Registry.Connect.SQLX().Get(&res, `select log_id, service, severity_level, info, details, occured, created
	from service_logs.logs where log_id = $1`,
		uid,
	); err != nil {

		return res, err

	}

	return res, nil

}

func (_ *sqlx_l) Store(p *protocols.Message) error {

	var res int = 0

	details, _ := json.MarshalIndent(p.Data, "", "	")

	if err := registry.Registry.Connect.SQLX().Get(&res,`insert into service_logs.logs (service, severity_level, info, details, occured)
		values ($1, $2, $3, $4, $5) returning log_id`,
			p.Owner,
			p.MessageType,
			p.Info,
			string(details),
			p.Time,
		); err != nil{

			return err
	}

	p.SetMonitorId(strconv.Itoa(res))

	return nil

}