package main

import (
	"controller/index"
	"controller/common"
	"app"
	"app/registry"
	"log"
	"flag"
	"net/http"
	"fmt"
	"time"
	//"app/registry"
	"runtime"
	"math"
	"os"
	"io/ioutil"
	_ "net/http/pprof"
	_ "lib/esi"
	_ "lib/scheduler"
	_ "lib/items"
	_ "lib/market"
	clib "lib/common"
	"controller/personal"
	"lib/items"
)

type Router interface {

	Route(a *app.App)
}

var (
	host         = flag.String("host", "", "host to listen to, leave blank to listen on any host")
	port         = flag.Uint("port", 80, "the port to listen on")
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

	registry.Registry.KVS.Put([]byte("started at"), time.Now())

	index.Route(a)
	common.Route(a)
	personal.Route(a)

	clib.Init()

	if err := items.GetDecryptorData(); err != nil {
		log.Fatal("Error getting decryptor data", err.Error())
	}

	fmt.Println("Starting HTTP server on port", *port)

	go func() {
		//pprof
		fmt.Println("Pprof listening on :8085")
		log.Println(http.ListenAndServe(":8085", nil))
	}()

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), "cert/cert.pem", "cert/key.pem", a.Router))

}