package notifies

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/gregdel/pushover"
)

var conf struct{
	AppToken string `yaml:"app_token"`
	Recipients []string `yaml:"recipients"`
	URL string `yaml:"url"`
	PushoverLevel int8 `yaml:"level"`
}

func init() {

	config := os.Getenv("CONFIG")

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/pushover.yml", config))

	if err != nil {

		log.Fatal(err.Error())
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {

		log.Fatal(err.Error())
	}

	pov = po{
		app : pushover.New(conf.AppToken),
	}

	for _, r := range conf.Recipients {
		pov.recipients = append(pov.recipients, pushover.NewRecipient(r))
	}

}
