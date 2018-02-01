package stats

import (
	"github.com/shirou/gopsutil/cpu"
	"errors"
	"fmt"
)

var times, timesC cpu.TimesStat

func CPUTimes() error {

	cts, err := cpu.Times(false)

	if err != nil {
		return errors.New(fmt.Sprintf("Read stat failed: %s", err.Error()))
	}

	timesC = times

	times = cts[0]

	return nil

}
