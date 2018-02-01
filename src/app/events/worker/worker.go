package worker

import (
	"app/events/event"
	"github.com/adjust/redismq"
	"log"
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
	"syscall"
)

var numWorkers = 2

func New(rmq func() *redismq.Queue, actioner Actioner, cosnumer Consumer) {

	l := listener{
		stop:         make(chan struct{}),
		stopped:      make(chan string),
		consumerName: cosnumer,
		queue:        rmq,
		writer:       actioner,
	}

	for i := 0; i < numWorkers; i++ {

		consumer, err := l.queue().AddConsumer(fmt.Sprintf("%s_%d", l.consumerName, i))

		if err != nil {

			fmt.Println(err.Error(), l.consumerName)

			return

		}

		go l.worker(consumer)
	}

	signalChan := make(chan os.Signal)

	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
	)

	<-signalChan
	// засылаем стоп
	for i := 0; i < numWorkers; i++ {

		l.stop <- struct{}{}
	}
	//и ждем подтверждения что можно "грохнуться"
	for i := 0; i < numWorkers; i++ {

		consumer := <-l.stopped

		log.Printf("Done: %s", consumer)
	}

	os.Exit(1)

}

func (l *listener) worker(consumer *redismq.Consumer) {

	fmt.Printf("Started worker %s\n", l.consumerName)

	packages := make(chan *redismq.Package)

	complete := make(chan struct{})

	go func() {

		for {
			// Unacked нафиг не нужны, т.к. с ними пока фиг его знает что делать, по идее можно отправить снова, ХЗ вобщем
			if consumer.HasUnacked() {

				log.Printf("Consumer (%s) Unacked", consumer.Name)

				if p, err := consumer.GetUnacked(); err == nil {

					p.Ack()
				}

				continue
			}

			p, err := consumer.Get()

			if err != nil {

				log.Printf("Consumer (%s) Get err: %s ", consumer.Name, err.Error())

				continue
			}

			packages <- p

			<-complete
		}
	}()


	for {

		select {

			case <-l.stop:

			l.stopped <- consumer.Name

				return

			case p := <-packages:

				event := event.Event{}

			if err := json.Unmarshal([]byte(p.Payload), &event); err == nil {

				if werr := l.writer.Action(event); werr != nil {

					//Потом разберемся что с ними делать, а если не ушла в очередь зафейленых - то и очень жаль
					p.Fail()

				}

			}

			p.Ack()

			complete <- struct{}{}
		}

	}

}
