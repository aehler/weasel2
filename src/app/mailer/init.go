package mailer

import (
	"app/events/manager"
	"log"
	"github.com/kataras/go-mailer"
	"fmt"
	"gopkg.in/yaml.v2"
	"app/events/event"
	"app"
	"github.com/flosch/pongo2"
)

type Creds struct {
	Sendmail struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		From     string `yaml:"from"`
	}
}

func Init(configData []byte) {

	rr := Creds{}

	if err := yaml.Unmarshal(configData, &rr); err != nil {

		log.Fatal(err.Error())
	}

	var err error

	q, err = manager.NewQueue(QueueName)
	if err != nil {
		log.Fatal(err)
	}

	cs1, err := manager.Subscribe(q.Name)
	if err != nil {
		log.Fatal(err)
	}

	go mail(cs1)

	sender = mailer.New(mailer.Config{
		Host:     rr.Sendmail.Host,
		Username: rr.Sendmail.Username,
		Password: rr.Sendmail.Password,
		FromAddr: rr.Sendmail.From,
		Port:     rr.Sendmail.Port,
		// Enable UseCommand to support sendmail unix command,
		// if this field is true then Host, Username, Password and Port are not required,
		// because these info already exists in your local sendmail configuration.
		//
		// Defaults to false.
		UseCommand: false,
	})

	fmt.Println(rr.Sendmail.From)

	log.Println("Mail app inited")

}

func mail(s chan manager.Msg) {

	var rendered string
	var err error

	for {

		p := <-s

		log.Println("Got message", p.Payload)

		if tn, ok := app.Templates[p.Payload.Object]; !ok {
			log.Println("Template not found:", p.Payload.Object)
			continue
		} else {
			rendered, err = tn.Execute(pongo2.Context(p.Payload.EventData))
			if err != nil {
				log.Println("Template not rendered:", err.Error())
				continue
			}
		}

		if err := sender.Send("Test Email", rendered, "example@example.com"); err != nil {

			log.Println("Error sending message", p.UUID, p.Payload, err.Error())

			continue
		}

		p.Ack()
	}

}

func SendQueued(object string, msg map[string]interface{}) {
	q.Pub(event.Event{
		Object: object,
		EventData: msg,
	})
}