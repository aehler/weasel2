package protocols

import "time"

const (
	N_TYPE_MESSAGE int8 = iota-1
	N_TYPE_WARNING
	N_TYPE_ERROR
	N_TYPE_PANIC
)

type Message struct {
	monitorId   string
	MessageType int8
	Owner       string
	Info        string
	Time		time.Time
	Data        interface{}
}

func (m *Message) SetMonitorId(id string) {
	m.monitorId = id
}

func (m *Message) MonitorId() string {
	return m.monitorId
}

func NewMessage(info, owner string, mt int8, data interface{}) *Message {

	return &Message{
		MessageType: mt,
		Owner: owner,
		Info: info,
		Data: data,
		Time: time.Now(),
	}

}