package manager

import (
	"app/events/event"
	"time"
	"github.com/juju/errors"
	"log"
	"fmt"
	"encoding/hex"
	"math/rand"
)

func NewQueue(name string) (*Queue, error) {

	if _, ok := queues[name]; ok {
		return nil, errors.New("[QM] Queue "+name+" already exists")
	}

	queues[name] = &Queue{
		Name: name,
		NumMsg: 0,
		new: make(chan struct{}, 1),
		msgs: &msgmx{
			msgs: make(map[string]*Msg),
		},
	}

	go queues[name].broadcast()

	log.Println("[QM] Queue "+ name+ " created")

	return queues[name], nil
}

func (q *Queue) broadcast() {

	fmt.Println("[QM] Broadcasting to", q.Name)

	for {

		<-q.new

		if len(q.subs) > 0 {

			q.msgs.mutex.Lock()
			for i, ms := range q.msgs.msgs {

				if ms.State != 0 {
					continue
				}

				for _, s := range q.subs {

					s <- *ms
					q.msgs.msgs[i].State = 1

				}

			}

			q.msgs.mutex.Unlock()
		}

	}

}

func (m *Msg) Ack() {

	q := queues[m.q]

	q.msgs.mutex.Lock()
	defer q.msgs.mutex.Unlock()

	if _, ok := q.msgs.msgs[m.UUID]; ok {
		delete(q.msgs.msgs, m.UUID)
	}

	q.NumMsg--
}

func Subscribe(name string) (chan Msg, error) {

	q, ok := queues[name]
	if !ok {
		return nil, errors.New("[QM] Queue "+name+" not found")
	}

	c := make(chan Msg)

	q.subs = append(q.subs, c)

	return c, nil
}

func (q *Queue) Pub(e event.Event) error {

	q.msgs.add(Msg{
		Payload: e,
		CreatedAt: time.Now(),
		State: 0,
		q: q.Name,
	})

	q.NumMsg++

	fmt.Println("[QM] Published to", q.Name, q.NumMsg)

	q.new <- struct {}{}

	return nil
}

func (m *msgmx) add(ms Msg) {

	u := make([]byte, 16)

	m.mutex.Lock()

	_, err := rand.Read(u)
	if err != nil {
		return
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	ms.UUID = hex.EncodeToString(u)

	m.msgs[ms.UUID] = &ms

	m.mutex.Unlock()
}