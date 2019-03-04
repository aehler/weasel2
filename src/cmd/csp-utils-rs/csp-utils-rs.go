package main

import (
	"log"
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"strings"
	exec2 "github.com/adlane/exec"

	"time"
	"strconv"
)

var conf struct {
	Csp struct {
		CprospPath string
	}
}

var (
	f            = flag.String("f", "", "Path to certificate file")
	inst         = flag.Bool("inst", false, "Install certificate")
	uninst       = flag.Bool("uninst", false, "Uninstall certificate")
	store        = flag.String("store", "root", "Storage")
	t            = flag.String("thumbprint", "", "thumbprint for deletion")
	config       = flag.String("config", "", "config")
	command		 = "certmgr"
	Ctx          exec2.ProcessContext
	done         chan struct{}
	errorc       chan struct{}
	attempt      = 1
	maxattempt   = 5
)

type reader struct {
}

func (*reader) OnData(b []byte) bool {

	fmt.Println("===")

	fmt.Print("MESSAGE:", string(b))

	if strings.Contains(string(b), "(o)OK, (c)Cancel") {

		fmt.Println("-= Sending o =-")

		Ctx.Send("o\n")

	}

	if strings.Contains(string(b), "Please choose index") {

		fmt.Println("-= Sending 2 =-")

		Ctx.Send("2\n")

	}

	if strings.Contains(string(b), "No certificate matching the criteria") {

		log.Println("Failed, do not continue")
		done <- struct{}{}

	}

	if strings.Contains(string(b), "[ErrorCode: 0x00000000]") {
		log.Println("Done")
		done <- struct{}{}
	}

	if strings.Contains(string(b), "[ErrorCode: 0x8010002c]") {
		log.Println("Failed, do not continue")
		done <- struct{}{}
	}

	return false
}

func (*reader) OnError(b []byte) bool {

	fmt.Print("ERROR:", string(b))

	if strings.Contains(string(b), "No certificate matching the criteria") {

		log.Println("Failed, do not continue")
		done <- struct{}{}

	}

	if strings.Contains(string(b), "[ErrorCode: 0x8010002c]") {
		log.Println("Failed, do not continue")
		done <- struct{}{}
	}

	if !strings.Contains(string(b), "WARNING") {
		errorc <- struct{}{}
	}
	return false
}


func (*reader) OnTimeout() {
	log.Println("Died on timeout")
	errorc <- struct{}{}
}

func checkRoot () {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	// output has trailing \n
	// need to remove the \n
	// otherwise it will cause error for strconv.Atoi
	// log.Println(output[:len(output)-1])

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Fatal(err)
	}

	if i != 0 {
		log.Fatal("This program must be run as root! (sudo)")
	}

}

func execute(command string, args ...string) {

ATTEMPTLOOP:
	for {

		log.Println("\n-----------------",attempt,"/",maxattempt,"-----------------")
		log.Println("Command:", command, args)

		Ctx = exec2.InteractiveExec(command, args...)

		r := reader{}
		go Ctx.Receive(&r, 10*time.Second)

		for {
			select {
			case <-done:
				os.Exit(0)
			case <-errorc:
				if attempt <= maxattempt {
					log.Println("Try again")
					attempt++
					Ctx.Cancel()
					Ctx.Stop()
					time.Sleep(1*time.Second)
					continue ATTEMPTLOOP
				}
				log.Println("Finished with error")
				os.Exit(1)
			}
		}
	}

	return

}

func main() {

	//checkRoot()

	done = make(chan struct{})
	errorc = make(chan struct{})

	flag.Parse()

	data, err := ioutil.ReadFile(fmt.Sprintf("%s", *config))

	if err != nil {

		log.Fatal(err.Error())
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		log.Fatal("Config parse error: ", err.Error())
	}

	if _, err := os.Stat(*f); err != nil && *inst {
		log.Fatal("Certificate file not found", *f)
	}

	if *uninst && *t == "" {
		log.Fatal("Must specify Thumbprint")
	}

	if *inst {
		execute(conf.Csp.CprospPath+command, "-inst", "-store", *store, "-f", *f)
	}

	if *uninst {
		execute(conf.Csp.CprospPath+command, "-delete", "-store", *store, "-thumbprint", *t)
	}

	log.Fatal("Unknown command, must specify -inst or -uninst, view --help for details")

	//go func() {
	//	//pprof
	//	fmt.Println("Pprof listening on :8085")
	//	log.Println(http.ListenAndServe(":8085", nil))
	//}()

}

