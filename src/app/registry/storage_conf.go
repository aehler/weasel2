package registry

import (
	"io/ioutil"
	"log"
	"fmt"
	"gopkg.in/yaml.v2"
)

type storage struct {
	Host string
	Dir string
}

func readStorageConf(config string) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/storage.yml", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	rr := map[string]*storage{}

	if err := yaml.Unmarshal(data, &rr); err != nil {

		log.Fatal(err.Error())
	}

	Registry.storage = rr
}
