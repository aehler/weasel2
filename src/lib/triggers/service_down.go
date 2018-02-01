package triggers

import (
	"lib/notifies"
	"github.com/akdcode/monitor/protocols"
	"lib/logs"
	"fmt"
	"log"
)

var serviceFailed = make(map[string]bool)

func serviceDown(t *triggerData) {

	t.mu.Lock()
	defer t.mu.Unlock()

	if(len(t.pid) == 0) {
		log.Println("Not enough data to trigger serviceDown")
		return
	}

	for _, serv := range t.pid[len(t.pid)-1] {

		if _, ok := serviceFailed[serv.Name]; !ok {

			serviceFailed[serv.Name] = false

		}

		//Up message
		if serviceFailed[serv.Name] == false && len(serv.PID) > 0 {

			p := protocols.NewMessage(
				fmt.Sprintf("Service %s UP", serv.Name),
				serv.Name,
				protocols.N_TYPE_MESSAGE,
				nil,
			)

			logs.Logs.Store(p)
			notifies.ForcePushover(p)
		}

		//Down message
		if serviceFailed[serv.Name] == true && len(serv.PID) == 0 {

			p := protocols.NewMessage(
				fmt.Sprintf("Service %s DOWN", serv.Name),
				serv.Name,
				protocols.N_TYPE_ERROR,
				serv.Error,
			)

			logs.Logs.Store(p)
			notifies.PushoverMQ(p)

		}

		if len(serv.PID) > 0 {

			serviceFailed[serv.Name] = true

			continue
		}

		serviceFailed[serv.Name] = false
	}

}