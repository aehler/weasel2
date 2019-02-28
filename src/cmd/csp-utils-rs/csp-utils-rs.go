package main

import (
	"log"
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"os/exec"
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
	config       = ""
	command		 = "crtmgr"
)

func main() {

	config = os.Getenv("CONFIG")

	data, err := ioutil.ReadFile(fmt.Sprintf("%s", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()

	if _, err := os.Stat(*f); err != nil {
		log.Fatal("Certificate file not found", *f)
	}

	if *inst {

		cmd := exec.Command(command, "-inst", "-store", *store, "-f", *f)
		//cmd := exec.Command("echo", "123")
		//cmd := exec.Command("openssl", "x509", "-in", "../cert/cert.pem", "-text")

		//stdin, err := cmd.StdinPipe()
		//if err != nil {
		//	log.Fatal(err.Error())
		//}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal("Stdout pipe init error ", err.Error())
		}


		defer func (){
			stdout.Close()
		}()

		err = cmd.Start()
		if err != nil {

			resp, err2 := ioutil.ReadAll(stdout)
			if err2 != nil {
				log.Println("Stdout read error ", err2.Error())
			}

			log.Println(string(resp))

			log.Fatal("Run error ", err.Error())
		}

		resp, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal("Read error ", err.Error())
		}

		log.Println(string(resp))

		if err := cmd.Wait(); err != nil {
			log.Fatal("Wait error ", err.Error())
		}

		return
	}

	if *uninst {
		return
	}

	log.Fatal("Unknown command, must specify -inst or -uninst, view --help for details")

	//go func() {
	//	//pprof
	//	fmt.Println("Pprof listening on :8085")
	//	log.Println(http.ListenAndServe(":8085", nil))
	//}()

}