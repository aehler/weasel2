package stats

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"fmt"
	"app/registry"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"errors"
	"log"
	"sync"
)

var ps = PidStatsM{
	Ps : make(map[string]PidStats),
}


type PidStats struct {
	PID   string
	State string
	Vsize uint64
	Rss   uint64
	Swap  uint64
	Utime uint64
	Stime uint64
	TimeP float64
}

type PidStatsM struct{
	Ps map[string]PidStats
	Mu sync.Mutex
}

func PS() map[string]PidStats {

	ps.Mu.Lock()
	defer ps.Mu.Unlock()

	pps := make(map[string]PidStats)

	for k, v := range ps.Ps {
		pps[k] = v
	}

	return pps

}

func (_ Metrics) PSStats() {

	percentage, err := cpu.Percent(0, true)
	if err != nil {
		log.Println(errors.New(fmt.Sprintf("Read stat failed: %s", err.Error())))
	}

	var total float64 = 0;

	for _, p := range percentage {
		total = total + p
	}

	ms, err := mem.VirtualMemory()
	if err == nil {
		registry.Registry.WsChan.PutAll("memStat", []uint64{ms.Available, ms.Used, ms.SwapCached})
		registry.Registry.WsChan.PutAll("cpuPercent", []float64{total, ms.UsedPercent})
	}

}

func (_ Metrics) TCPSockets() {

	if stats, err := linuxproc.ReadNetTCPSockets("/proc/net/tcp", linuxproc.NetIPv4Decoder); err == nil {

		registry.Registry.WsChan.PutAll("openTCPSockets", len(stats.Sockets))

	} else {

		log.Println(err.Error())

	}

}

func (_ Metrics) PSStatsPID() {

	ps.Mu.Lock()
	Pids.Mu.Lock()

	defer ps.Mu.Unlock()
	defer Pids.Mu.Unlock()

	busy := times.User + times.System + times.Nice + times.Iowait + times.Irq +
		times.Softirq + times.Steal + times.Guest + times.GuestNice + times.Stolen
	busyC := timesC.User + timesC.System + timesC.Nice + timesC.Iowait + timesC.Irq +
		timesC.Softirq + timesC.Steal + timesC.Guest + timesC.GuestNice + timesC.Stolen

	for i, p := range Pids.P {

		for _, pid := range p.PID {

			if stats, err := linuxproc.ReadProcessStat(fmt.Sprintf("/proc/%s/stat", pid)); err == nil {

				statm, err := linuxproc.ReadProcessStatus(fmt.Sprintf("/proc/%s/status", pid));

				if err != nil {

					log.Println(err.Error())

					Pids.P[i].Error = err.Error()

					continue

				}

				var utime uint64

				if psc, ok := ps.Ps[pid]; ok {
					utime = psc.Utime
				}

				np := PidStats{
					State: stats.State,
					Vsize: statm.VmData + statm.VmExe + statm.VmStk,
					Rss  : uint64(statm.VmRSS),
					Swap : statm.VmSwap,
					Utime: stats.Utime,
					Stime: stats.Stime,
					TimeP: float64(float64(stats.Utime - utime) / ( ((times.Idle + busy) - (timesC.Idle + busyC)) * 100 )),
				}

				ps.Ps[pid] = np

			} else {

				log.Println(err.Error())

				Pids.P[i].Error = err.Error()

			}

		}

	}

	ex := []string{}

	//clean ps
	for rp, _ := range ps.Ps {

		for _, p := range Pids.P {

			for _, pid := range p.PID {

				if pid == rp {

					ex = append(ex, rp)

				}
			}
		}
	}

	ML:
	for p, _ := range ps.Ps {

		for _, k := range ex {

			if k == p {

				continue ML
			}

		}

		delete(ps.Ps, p)

	}

	var master = make(map[string]PidStats)

	for _, p := range Pids.P {

		pp := PidStats{}

		for _, pid := range p.PID {

			if pss, ok := ps.Ps[pid]; ok {

				pp.TimeP = pss.TimeP + pp.TimeP
				pp.Vsize = pss.Vsize + pp.Vsize
				pp.Rss   = pss.Rss + pp.Rss
				pp.Swap  = pss.Swap + pp.Swap

			}

		}

		master[p.Name] = pp

	}

	registry.Registry.WsChan.PutAll("ps", master)

}