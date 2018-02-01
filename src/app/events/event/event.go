package event

import (
	"github.com/akdcode/srm-services/auth"
)

type EventType string

type Event struct {
	Object string
	ObjectId uint
	EventType EventType
	EventData map[string]interface{}
	User *auth.User
}
