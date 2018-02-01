package registry

import (
	"io/ioutil"
	"log"
	"fmt"
	"gopkg.in/yaml.v2"
)

type refConf struct {
	RefType string `yaml:"type"`
	Fields map[string]refField `yaml:"fields"`
}

type refField struct {
	Type string  `yaml:"type"`
	Name string  `yaml:"name"`
	Label string `yaml:"label"`
	Ord uint    `yaml:"ord"`
}

func readRefConf(config string) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/classifiers.yml", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	rr := map[string]*refConf{}

	if err := yaml.Unmarshal(data, &rr); err != nil {

		log.Fatal(err.Error())
	}

	Registry.ReferenceConf = rr
}

