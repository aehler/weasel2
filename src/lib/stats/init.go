package stats

import (
	"log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const MONITOR_TICKS_SECONDS time.Duration = 5 * time.Second

type SysStats struct {
	MemTotal uint64
	CPU []cpu.InfoStat
}

var GenStats SysStats

type Metrics struct{}

func init() {

	stat, err := cpu.Info()
	if err != nil {
		log.Fatal("stat read fail")
	}

	ms, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal("stat read fail")
	}

	GenStats = SysStats{
		MemTotal : ms.Total,
		CPU : stat,
	}

	Pids = PidsM{
		P : []Pid{},
	}

}