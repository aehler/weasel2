package triggers

import (
	"log"
	"time"
	"lib/stats"
)

func (t *triggerData) run() {

	log.Println("Running triggers after one second...")

	time.Sleep(time.Second)

	ps := stats.PS()
	pids := stats.GetPidsM()

	time.Sleep(stats.MONITOR_TICKS_SECONDS)

	for {

		ps = stats.PS()
		pids = stats.GetPidsM()

		t.addPids(pids)

		t.addStats(ps)

		for _, c := range checks {

			c(t)

		}

		time.Sleep(stats.MONITOR_TICKS_SECONDS)
	}
}