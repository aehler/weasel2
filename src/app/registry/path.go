package registry

import (
	"log"
	"gopkg.in/yaml.v2"
)

type Path struct {
	Templates string
	Static string
	HTTPStatic string `yaml:"HTTPStatic"`
}

func ReadPathConf(data []byte) *Path {

	rr := map[string]*Path{}

	if err := yaml.Unmarshal(data, &rr); err != nil {

		log.Fatal(err.Error())
	}

	return rr["path"]
}
