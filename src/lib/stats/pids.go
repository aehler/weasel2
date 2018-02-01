package stats

import (
	"app/registry"
	"os/exec"
	"strings"
	"github.com/shirou/gopsutil/net"
	"strconv"
	"log"
	"sync"
)

var Pids PidsM

type Pid struct{
	PID   []string
	Name  string
	Ping  bool
	Error string
	Conns uint
}

type PidsM struct{
	P []Pid
	Mu sync.Mutex
}

func GetPidsM() []Pid {

	Pids.Mu.Lock()
	defer Pids.Mu.Unlock()

	res := make([]Pid, len(Pids.P))
	copy(res, Pids.P)

	return res
}

func UpdatePids() error {

	Pids.Mu.Lock()
	defer Pids.Mu.Unlock()

	Pids.P = []Pid{}

	for _, p := range registry.Registry.Monitor {

		if s, err := exec.Command("pgrep", "^"+p).Output(); err != nil {

			Pids.P = append(Pids.P, Pid{
				PID:   []string{},
				Name:  p,
				Error: err.Error(),
			})

		} else {

			ss := strings.Trim(string(s), "\n")

			allPids := strings.Split(ss, "\n")

			ptc := 0

			for _, pfn := range allPids {

				if pi, err := strconv.Atoi(pfn); err == nil {

					if netstat, err := net.ConnectionsPid("tcp", int32(pi)); err == nil {

						ptc = ptc + len(netstat)

					} else {

						log.Println(err.Error())

					}

				}

			}

			Pids.P = append(Pids.P, Pid{
				PID:   allPids,
				Name:  p,
				Error: "",
				Conns: uint(ptc),
			})

		}
	}

	return nil

}

func (_ Metrics) SendPidInfo() {
	registry.Registry.WsChan.PutAll("pids", GetPidsM())
}