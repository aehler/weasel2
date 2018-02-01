package triggers

import (
	"os/exec"
	"app/registry"
	"time"
	"lib/logs"
	"github.com/akdcode/monitor/protocols"
	"fmt"
	"lib/notifies"
)

func RestartByFailure(p *protocols.Message) error {

	switch p.Info {

	case
	//"Too many GC calls",
	"socket: too many open files":

		var p *protocols.Message

		if err := RestartService(p.Owner); err != nil {

			p = protocols.NewMessage(
				"Autorestart by monitor failed",
				p.Owner,
				protocols.N_TYPE_PANIC,
				p.Info,
			)

		} else {

			p = protocols.NewMessage(
				"Autorestart by monitor",
				p.Owner,
				protocols.N_TYPE_MESSAGE,
				p.Info,
			)

		}

		logs.Logs.Store(p)
		notifies.ForcePushover(p)

	}

	return nil
}

func RestartService(s string) (Err error) {

	defer func(){

		msg := fmt.Sprintf("Service %s restarted by monitor", s)
		severity := protocols.N_TYPE_MESSAGE
		error := ""

		if Err != nil {
			msg = fmt.Sprintf("Failed to restart %s by monitor", s)
			severity = protocols.N_TYPE_WARNING
			error = Err.Error()
		}

		p := protocols.NewMessage(
			msg,
			s,
			severity,
			error,
		)

		logs.Logs.Store(p)
		notifies.WS(p)

	}()

	switch s {

	case "monitor":
		return

	default:

		if _, Err = exec.Command("systemctl", "stop", registry.Registry.SysctlNames[s]).Output(); Err == nil {

			time.Sleep(time.Second * 2)

			if _, Err = exec.Command("systemctl", "start", registry.Registry.SysctlNames[s]).Output(); Err == nil {

				return

			} else {

				return

			}

		} else {

			return

		}

	}

}