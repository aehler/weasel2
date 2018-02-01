package notifies

import (
	"fmt"
	"time"
	"github.com/gregdel/pushover"
	"github.com/akdcode/monitor/protocols"
	"log"
)

var pov po

type po struct {
	app        *pushover.Pushover
	recipients []*pushover.Recipient
}

func (p po) AppToken() string {
	return conf.AppToken
}

func ForcePushover(p *protocols.Message) {

	pov.send(p)

}

func PushoverMQ(p *protocols.Message) error {

	if conf.PushoverLevel > p.MessageType {
		return nil
	}

	pov.send(p)

	return nil
}

func StatusPushOver() []*pushover.RecipientDetails {

	res := []*pushover.RecipientDetails{}

	for _, r := range pov.recipients {

		log.Println(r)

		recipientDetails, err := pov.app.GetRecipientDetails(r)
		if err != nil {
			log.Println(err.Error())
			continue;
		}

		res = append(res, recipientDetails)
	}

	return res
}

func (p po) send(msg *protocols.Message) {

	message := &pushover.Message{
		Message:     msg.Info,
		Title:       fmt.Sprintf("Сообщение сервиса %s", msg.Owner),
		Priority:    int(msg.MessageType),
		URL:         fmt.Sprintf("%s/message/%d/", conf.URL, msg.MonitorId()),
		URLTitle:    "Просмотр полной информации об ошибке",
		Timestamp:   time.Now().Unix(),
		Retry:       300 * time.Second,
		Expire:      time.Hour * 24,
		Sound:       pushover.SoundCosmic,
	}

	// Send the message to each recipient
	for _, r := range p.recipients{

		response, err := p.app.SendMessage(message, r)
		if err != nil {
			log.Println("Error sending message", err.Error())
		}

		log.Println(response)

	}

}