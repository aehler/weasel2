package events

import (
	"github.com/akdcode/srm-lib/business_events/event"
	"github.com/adjust/redismq"
	"encoding/json"
	"fmt"
)

func New(rmq func() *redismq.Queue) *Firer {

	fmt.Printf("starting firer %s...", rmq().Name)
	fmt.Println("")

	return &Firer {
		queue : rmq,
	}
}

type Firer struct {
	queue func() *redismq.Queue
}

func (f *Firer) Fire (e event.Event) error {

	b, err := json.Marshal(e)

	if err != nil {

		fmt.Println(err)

		return err
	}

	f.queue().Put(string(b))

	return nil

}
