package triggers

import (
	"lib/stats"
	"lib/logs"
	"fmt"
	"github.com/akdcode/monitor/protocols"
	"runtime/debug"
)

var td *triggerData
var checks []trigger = []trigger{}
var gs *stats.SysStats

func init() {

	defer func() {
		if r := recover(); r != nil {

			p := protocols.NewMessage(
				"Monitor panic recovered",
				"monitor",
				protocols.N_TYPE_ERROR,
				[]string{fmt.Sprintf("%v", r), string(debug.Stack())},
			)

			logs.Logs.Store(p)

		}
	}()

	td = &triggerData{}

	gs = &stats.GenStats

	checks = append(checks,
		serviceDown,
		cpuLimit,
		memLimit,
		)

	go td.run()

}