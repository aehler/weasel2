package session

import (
	"github.com/bradfitz/gomemcache/memcache"
	"encoding/json"
	"fmt"
)

type SessionStorage struct {
	mcClient *memcache.Client
}

func Init() (*SessionStorage) {

	s := &SessionStorage{
		mcClient : memcache.New("127.0.0.1:11211"),
	}

	return s
}

func (s *SessionStorage) Add(ssid string, user interface {}) error {

	fmt.Println("Setting new", ssid)

	b, err := json.Marshal(user)

	if err != nil {

		return err
	}

	err = s.mcClient.Add(&memcache.Item{
		Key : ssid,
		Value : b,
		Flags : 0,
		Expiration : 3600,
	})
	if err != nil {

		return err
	}

	return nil
}

func (s *SessionStorage) Get(ssid string) ([]byte, error) {

	i, err := s.mcClient.Get(ssid)

	if err != nil {

		return nil, err

	}

	s.mcClient.Touch(ssid, 3600)

	return i.Value, nil
}

func (s *SessionStorage) Kill(ssid string) {

	err := s.mcClient.Delete(ssid)

	if err != nil {

		fmt.Println(err)
	}

}
