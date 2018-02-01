package app

import (
	"encoding/json"
	"expvar"
	"runtime"
	"strconv"
	"time"
)

var (
	requestCount = expvar.NewMap("req_count")
	requestTime  = expvar.NewMap("req_time")
)

func timerStart(urlPath string) func() {

	start := time.Now()

	return func() {

		requestCount.Add(urlPath, 1)
		requestTime.Add(urlPath, int64(time.Now().Sub(start)))
	}
}

type MemStats struct {
	// General statistics.
	Alloc      uint64 // bytes allocated and still in use
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups    uint64 // number of pointer lookups
	Mallocs    uint64 // number of mallocs
	Frees      uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and still in use
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  uint64 // bytes used by stack allocator
	StackSys    uint64
	MSpanInuse  uint64 // mspan structures
	MSpanSys    uint64
	MCacheInuse uint64 // mcache structures
	MCacheSys   uint64
	BuckHashSys uint64 // profiling bucket hash table
	GCSys       uint64 // GC metadata
	OtherSys    uint64 // other system allocations

	// Garbage collector statistics.
	NextGC       uint64 // next collection will happen when HeapAlloc ≥ this amount
	LastGC       uint64 // end time of last collection (nanoseconds since 1970)
	PauseTotalNs uint64
	//	PauseNs      [256]uint64 // circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]
	//	PauseEnd     [256]uint64 //
	NumGC uint32
}

type metrics struct {
	Version                string
	NumGoroutine           int
	RequestsTotalCount     uint64
	RequestsTotalTime      uint64
	RequestCountByLocation map[string]uint64
	RequestTimeByLocation  map[string]uint64
	MemStats               MemStats
}

func metricsHandler(c *Context) {

	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	m := metrics{
		Version:                runtime.Version(),
		NumGoroutine:           runtime.NumGoroutine(),
		RequestCountByLocation: make(map[string]uint64),
		RequestTimeByLocation:  make(map[string]uint64),
		MemStats: MemStats{
			Alloc:        memStats.Alloc,
			TotalAlloc:   memStats.TotalAlloc,
			Sys:          memStats.Sys,
			Lookups:      memStats.Lookups,
			Mallocs:      memStats.Mallocs,
			Frees:        memStats.Frees,
			HeapAlloc:    memStats.Frees,
			HeapSys:      memStats.HeapSys,
			HeapIdle:     memStats.HeapIdle,
			HeapInuse:    memStats.HeapInuse,
			HeapReleased: memStats.HeapReleased,
			HeapObjects:  memStats.HeapObjects,
			StackInuse:   memStats.StackInuse,
			StackSys:     memStats.StackSys,
			MSpanInuse:   memStats.MSpanInuse,
			MSpanSys:     memStats.MSpanSys,
			MCacheInuse:  memStats.MCacheInuse,
			MCacheSys:    memStats.MCacheSys,
			BuckHashSys:  memStats.BuckHashSys,
			GCSys:        memStats.GCSys,
			OtherSys:     memStats.OtherSys,
			NextGC:       memStats.NextGC,
			LastGC:       memStats.LastGC,
			PauseTotalNs: memStats.PauseTotalNs,
			NumGC:        memStats.NumGC,
		},
	}

	requestCount.Do(func(k expvar.KeyValue) {

		m.RequestCountByLocation[k.Key] = metricToUnt64(k.Value.String())

		m.RequestsTotalCount += m.RequestCountByLocation[k.Key]
	})

	requestTime.Do(func(k expvar.KeyValue) {

		m.RequestTimeByLocation[k.Key] = metricToUnt64(k.Value.String())

		m.RequestsTotalTime += m.RequestTimeByLocation[k.Key]
	})

	json.NewEncoder(c).Encode(m)
}

func metricToUnt64(m string) uint64 {

	v, _ := strconv.ParseUint(m, 10, 0)

	return v
}
