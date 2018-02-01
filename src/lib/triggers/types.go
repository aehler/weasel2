package triggers

import (
	"sync"
	"lib/stats"
)

const SHIFT_TICKS int = 6
const LOAD_THRESHOLD_PERCENT float64 = 75
const MEM_SERVICE_THRESHOLD_PERCENT float64 = 50

type triggerData struct {
	pidStats []map[string]stats.PidStats
	pid [][]stats.Pid
	mu sync.Mutex
}

type trigger func(t *triggerData)

func (t *triggerData) addStats(s map[string]stats.PidStats) {

	t.mu.Lock()
	defer t.mu.Unlock()

	//nm := make(map[string]stats.PidStats)
	//
	//for k, v := range s {
	//	nm[k] = v
	//}

	t.pidStats = append(t.pidStats, s)

	if len(t.pidStats) > SHIFT_TICKS {

		t.pidStats = t.pidStats[len(t.pidStats)-SHIFT_TICKS:len(t.pidStats)]

	}

}

func (t *triggerData) addPids(s []stats.Pid) {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.pid = append(t.pid, s)

	if len(t.pid) > SHIFT_TICKS {

		t.pid = t.pid[len(t.pid)-SHIFT_TICKS:len(t.pid)]

	}

}