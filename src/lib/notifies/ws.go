package notifies

import (
	"github.com/akdcode/monitor/protocols"
	"app/registry"
)

func WS(p *protocols.Message) error {

	registry.Registry.WsChan.Put("serviceFailure", map[string]interface{}{
		"service" : p.Owner,
		"severity": p.MessageType,
		"id": p.MonitorId(),
		"name": p.Info,
	})

	return nil
}