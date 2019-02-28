package mailer

import (
	"app/events/manager"
	"github.com/kataras/go-mailer"
)

const QueueName string = "MAIL_Q"
var q *manager.Queue
var sender *mailer.Mailer