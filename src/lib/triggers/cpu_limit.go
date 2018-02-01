package triggers

import (
	"time"
	"fmt"
	"github.com/akdcode/monitor/protocols"
	"lib/notifies"
	"lib/logs"
)

var cpuLimitSent time.Time = time.Now()

func cpuLimit(t *triggerData) {

	t.mu.Lock()
	defer t.mu.Unlock()

	//Send pushover once every 10 minutes
	if time.Now().Before(cpuLimitSent.Add(time.Minute * 10)) {

		return

	}

	cpuLimits := map[string][]bool{}
	cpuLimitsAvg := map[string][]float64{}

	for _, iter := range t.pidStats {

		for name, stats := range iter {

			if _, ok := cpuLimits[name]; !ok {

				cpuLimits[name] = []bool{}

			}

			cpuLimitsAvg[name] = append(cpuLimitsAvg[name], stats.TimeP*100)

			if stats.TimeP*100 > LOAD_THRESHOLD_PERCENT {

				cpuLimits[name] = append(cpuLimits[name], true)

			} else {

				cpuLimits[name] = append(cpuLimits[name], false)

			}

		}

	}

	NAMELOOP:
	for name, d := range cpuLimits {

		if len(d) < SHIFT_TICKS {
			continue NAMELOOP
		}

		for _, dd := range d {

			if !dd {

				continue NAMELOOP

			}

		}

		nm := name

		SNL:
		for _, p := range t.pid[0] {

			for _, pp := range p.PID {

				if pp == name {

					nm = p.Name

					break SNL

				}

			}

		}

		var avg float64

		for _, a := range cpuLimitsAvg[name] {
			avg = avg + a
		}

		p := protocols.NewMessage(
			fmt.Sprintf("Service %s consuming more then %.0f%% cpu (%.2f%%)", nm, LOAD_THRESHOLD_PERCENT, (avg/float64(len(cpuLimitsAvg[name]))) ),
			nm,
			protocols.N_TYPE_ERROR,
			nil,
		)

		logs.Logs.Store(p)
		notifies.PushoverMQ(p)
		notifies.WS(p)

		cpuLimitSent = time.Now()
	}

}