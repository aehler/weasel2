package manager

import (
	"time"
	"sync"
	"app/events/event"
)

type Queue struct {
	Name string
	NumMsg uint64
	msgs *msgmx
	subs []chan Msg
	new chan struct{}
}

type msgmx struct {
	msgs map[string]*Msg
	mutex sync.Mutex
}

type Msg struct {
	Payload event.Event
	CreatedAt time.Time
	UUID string
	State uint8
	q string
}

var queues map[string]*Queue
