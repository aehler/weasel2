package registry

import (
	"fmt"
	"log"
	"github.com/go-redis/redis"
	"github.com/adjust/redismq"
	"time"
	"encoding/json"
	"github.com/akdcode/monitor/protocols"
)

type Redis struct {
	client *redis.Client
	MQ *redismq.Queue
	Listeners []func(p *protocols.Message) error
}

type payload struct {
	CreatedAt time.Time
	Data string
}

func (r *registry) newRedisClient(mqs *Rediscreds) {

	if mqs == nil {
		log.Println("Redis conf not found, some functions will be disabled")
		return
	}

	var queue *redismq.Queue

	if mqs.MQName != "" {

		queue = redismq.CreateQueue(mqs.Host, mqs.Port, mqs.Password, int64(mqs.DB), mqs.MQName)

	} else {

		log.Println("MQ config not found, error messages from services will not be served")

	}


	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", mqs.Host, mqs.Port),
		Password: "", // no password set
		DB:       mqs.DB,
	})

	pong, err := client.Ping().Result()

	if err != nil {
		log.Fatal("Connect to redis failed", err.Error())
	}

	fmt.Println("Redis: ", pong)

	r.Redis = &Redis{
		client: client,
		MQ: queue,
	}

	go r.Redis.readMQ()

	return

}

func (r *Redis) Store(key string, v interface{}) {

	d, err := json.Marshal(v)

	if err != nil {
		log.Println("redis store failed,", err.Error())
	}

	s, err := json.Marshal(payload{
		CreatedAt: time.Now(),
		Data: string(d),
	})

	r.client.LPush(key, string(s))

}

func (r *Redis) Get(key string, start, stop int) []string {

	list, err := r.client.LRange(key, 0, -1).Result()

	if err != nil {
		log.Println("redis lrange failed,", err.Error())

		return []string{}
	}

	return list
}

func (r *Redis) Trim(key string, len int64) {

	r.client.LTrim(key, 0, len).Result()

}

func (r *Redis) readMQ() {

	if r.MQ == nil {

		return
	}

	consumer, err := r.MQ.AddConsumer("critical-notify consumer")

	if err != nil {
		log.Fatal("RedisMQ consumer", err.Error())
	}

	log.Println("Added consumer")

	consumer.ResetWorking()

	for {

		p, err := consumer.Get()

		if err != nil {
			log.Println(err)
			continue
		}

		msg := &protocols.Message{}

		if err := json.Unmarshal([]byte(p.Payload), msg); err != nil {

			log.Println("RedisMQ unmarshal error", err.Error())

			p.Fail()

			continue
		}

		for _, l := range r.Listeners {

			if err := l(msg); err != nil {

				log.Println("RedisMQ listener error", err.Error())

			}

		}

		err = p.Ack()
	}

}

func (r *Redis) AddListener(f func(p *protocols.Message) error) {

	r.Listeners = append(r.Listeners, f)

}