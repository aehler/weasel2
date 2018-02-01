package main

import (
	"controller/index"
	"app"
	"app/registry"
	"log"
	"flag"
	"net/http"
	"fmt"
	"time"
	_ "lib/stats"
	"lib/stats"
	_ "lib/notifies"
	"lib/notifies"
	_"lib/triggers"
	//"app/registry"
	"golang.org/x/net/websocket"
	"sync"
	"runtime"
	"math"
	"os"
	"lib/logs"
	"lib/triggers"
	"io/ioutil"
	"reflect"
	_ "net/http/pprof"
)

type Router interface {

	Route(a *app.App)
}

var (
	host         = flag.String("host", "", "host to listen to, leave blank to listen on any host")
	port         = flag.Uint("port", 80, "the port to listen on")
	ws_port      = flag.Uint("ws_port", 3000, "the port to listen for websocket")
	withBinData  = flag.Bool("withbinstatic", false,"Use binary http static and templates")
	config       = ""
)

func main() {

	numCPUs := runtime.NumCPU()
	gmp := int(math.Ceil(float64(numCPUs/2)))
	runtime.GOMAXPROCS(gmp)

	config = os.Getenv("CONFIG")

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/config.yml", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	log.Printf("Running on %d CPUs, with %d max procs.\n", numCPUs, gmp)

	flag.Parse()

	a := app.New(data, config, *withBinData)

	logs.Init()

	if registry.Registry.Redis != nil {

		registry.Registry.Redis.AddListener(logs.Logs.Store)
		registry.Registry.Redis.AddListener(notifies.PushoverMQ)
		registry.Registry.Redis.AddListener(notifies.WS)
		registry.Registry.Redis.AddListener(triggers.RestartByFailure)

	}

	registry.Registry.KVS.Put([]byte("ws_port"), *ws_port)

	index.Route(a)

	fmt.Println("Starting HTTP server on port", *port)

	go func() {
		//pprof
		fmt.Println("Pprof listening on :8085")
		log.Println(http.ListenAndServe(":8085", nil))
	}()

	go monitor()

	go ws()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), a.Router))

}

func ws() {

	http.Handle("/", websocket.Handler(index.WsHandler))

	err := http.ListenAndServe(fmt.Sprintf(":%d", *ws_port), nil)

	if err != nil {

		panic("ListenAndServe: " + err.Error())

	}

}

func monitor() {

	metrics := []reflect.Value{}

	var syn = sync.Mutex{}

	s := stats.Metrics{}

	for _, method := range registry.Registry.Config.Metrics {

		m := reflect.ValueOf(s).MethodByName(method)

		if m.IsValid() {

			metrics = append(metrics, m)

			log.Println("Using", method)
		}

	}

	for {

		syn.Lock()
		stats.UpdatePids()
		stats.CPUTimes()
		syn.Unlock()

		for _, m := range metrics {

			go m.Call([]reflect.Value{})

		}

		//go s.SendPidInfo()
		//go s.PSStats()
		//go s.PSStatsPID()
		//go s.TCPSockets()

		//log.Printf("Stats read in %v\n", time.Since(t))

		time.Sleep(stats.MONITOR_TICKS_SECONDS)
	}

}