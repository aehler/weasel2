package worker

import (
	"github.com/adjust/redismq"
	"github.com/jmoiron/sqlx"
	"app/events/event"
)

const (
	BusinessEvents Consumer = "business-events"
	AroundProcedureEvents Consumer = "around-procedure-events"
	CommissionSessionEvents Consumer = "commission-session-events"
)

type listener struct {
	stop         chan struct{}
	stopped      chan string
	consumerName Consumer
	queue        func() *redismq.Queue
	db           func() *sqlx.DB
	writer       Actioner
}

type Consumer string

type Actioner interface {
	Action(e event.Event) error
}
