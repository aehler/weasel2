package esi

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"app/jsonRPC"
)

var conf struct{
	URL string `yaml:"esi"`
}

type Client struct {
	rpc       jsonRPC.Client
}

var MC Client

func init() {

	config := os.Getenv("CONFIG")

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/config.yml", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {

		log.Fatal(err.Error())
	}

	MC.rpc = jsonRPC.New("Market", conf.URL)
}
