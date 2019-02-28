package manager

import (
	"time"
	"fmt"
)

func mgc() {

	for {

		time.Sleep(600 * time.Second)

		fmt.Println("[QM] Cleanining outdated messages")

		for _, q := range queues {
			q.msgs.mutex.Lock()

			oldnum := q.NumMsg

			for _, msg := range q.msgs.msgs {

				if msg.State != 0 {
					delete(q.msgs.msgs, msg.UUID)
					q.NumMsg--
				}

			}
			q.msgs.mutex.Unlock()

			fmt.Sprintf("[QM] Cleaned up %d messages from queue %s, current message count %d \n", oldnum-q.NumMsg, q.Name, q.NumMsg)
		}
	}
}