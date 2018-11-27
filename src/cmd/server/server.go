package main

import (
	"app"
	"log"
	"flag"
	"net/http"
	"fmt"
	"runtime"
	"math"
	"os"
	"io/ioutil"
	_ "net/http/pprof"
	_ "lib/scheduler"
	"controller"
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

	controller.Route(a)

	fmt.Println("Starting HTTP server on port host", *host, ":", *port)

	go func() {
		//pprof
		fmt.Println("Pprof listening on :8085")
		log.Println(http.ListenAndServe(":8085", nil))
	}()

	log.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), a.Router))

	//log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), "cert/cert.pem", "cert/key.pem", a.Router))

}