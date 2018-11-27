package session

import (
	"github.com/bradfitz/gomemcache/memcache"
	"encoding/json"
	"fmt"
	"time"
)

type SessionStorage struct {
	mcClient *memcache.Client
	SStimeout int32
}

func Init(mg string, sst int32) (*SessionStorage) {

	s := &SessionStorage{
		mcClient : memcache.New(mg),
		SStimeout: sst,
	}

	if err := s.Upsert("Runtime_start", time.Now().String()); err != nil {
		panic(err.Error())
	}

	return s
}

func (s *SessionStorage) Cache(key string, value []byte, ttl time.Duration) {

	err := s.mcClient.Add(&memcache.Item{
		Key : key,
		Value : value,
		Flags : 0,
		Expiration : int32(ttl.Seconds()),
	})
	if err != nil {

		err = s.mcClient.Replace(&memcache.Item{
			Key : key,
			Value : value,
			Flags : 0,
			Expiration : int32(ttl.Seconds()),
		})

	}

}

func (s *SessionStorage) Upsert(key string, value interface{}) error {

	if _, err := s.mcClient.Get(key); err != nil {

		if err.Error() == "memcache: cache miss" {

			return s.Add(key, value)

		} else {

			return err

		}

	} else {

		return s.Replace(key, value)

	}

}

func (s *SessionStorage) Replace(key string, value interface{}) error {

	b, err := json.Marshal(value)

	if err != nil {

		return err
	}

	err = s.mcClient.Replace(&memcache.Item{
		Key : key,
		Value : b,
		Flags : 0,
		Expiration : s.SStimeout,
	})
	if err != nil {

		return err
	}

	return nil
}

func (s *SessionStorage) Add(ssid string, user interface {}) error {

	b, err := json.Marshal(user)

	if err != nil {

		return err
	}

	err = s.mcClient.Add(&memcache.Item{
		Key : ssid,
		Value : b,
		Flags : 0,
		Expiration : s.SStimeout,
	})
	if err != nil {

		return err
	}

	return nil
}

func (s *SessionStorage) GetNoTouch(ssid string) ([]byte, error) {

	i, err := s.mcClient.Get(ssid)

	if err != nil {

		return nil, err

	}

	return i.Value, nil
}

func (s *SessionStorage) Unmarshal(key string, res interface{}) error {

	b, err := s.GetNoTouch(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, res); err != nil {
		return err
	}

	return nil
}

func (s *SessionStorage) Get(ssid string) ([]byte, error) {

	i, err := s.mcClient.Get(ssid)

	if err != nil {

		return nil, err

	}

	s.mcClient.Touch(ssid, s.SStimeout)

	return i.Value, nil
}

func (s *SessionStorage) Kill(ssid string) {

	err := s.mcClient.Delete(ssid)

	if err != nil {

		fmt.Println(err)
	}

}
