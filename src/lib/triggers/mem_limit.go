package triggers

import (
	"time"
	"fmt"
	"github.com/akdcode/monitor/protocols"
	"lib/logs"
	"lib/notifies"
)

var memLimitSent time.Time

func memLimit(t *triggerData) {

	t.mu.Lock()
	defer t.mu.Unlock()

	//Send pushover once every 10 minutes
	if time.Now().Before(memLimitSent.Add(time.Minute * 10)) {

		return

	}

	cpuLimits := map[string][]bool{}
	cpuLimitsAvg := map[string][]float64{}

	for _, iter := range t.pidStats {

		for name, stats := range iter {

			if _, ok := cpuLimits[name]; !ok {

				cpuLimits[name] = []bool{}

			}

			cpuLimitsAvg[name] = append(cpuLimitsAvg[name], float64(stats.Vsize + stats.Rss)/1024/1024)

			if (float64(stats.Vsize + stats.Rss) / float64(gs.MemTotal))*100 > MEM_SERVICE_THRESHOLD_PERCENT {

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
			fmt.Sprintf("Service %s consuming more then %.0f%% total memory (%.0fMB of %dMB total)", nm, MEM_SERVICE_THRESHOLD_PERCENT, (avg/float64(len(cpuLimitsAvg[name]))), gs.MemTotal/1024/1024),
			nm,
			protocols.N_TYPE_ERROR,
			nil,
		)

		logs.Logs.Store(p)
		notifies.PushoverMQ(p)
		notifies.WS(p)

		RestartService(p.Owner)

		memLimitSent = time.Now()
	}

}