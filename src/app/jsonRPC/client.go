package jsonRPC

import (
	"fmt"
	"net/http"
	"errors"
	"encoding/json"
	"bytes"
	"time"
	"app/registry"
)

func New(name, address string) Client {

	fmt.Println("Starting client for", name)

	c := &client{
		name:    name,
		address: address,
	}

	return c
}

func (c *client) Send(params Params, result interface{}) error {

	var (
		sendErr  error
	)

	sec := make(chan error)

	go func() {

		//try from cache
		if bc, err := registry.Registry.Session.Get(c.address + params.Url()); err == nil {

			fmt.Println("Got from cache")

			if err := json.Unmarshal(bc, &result); err == nil {

				sec <- nil

				return

			} else {

				fmt.Println("Cache unmarshal error, get from esi", err.Error())

			}

		}

		client := &http.Client{}

		reqBody := bytes.NewBuffer([]byte{})

		req, err := http.NewRequest("GET", c.address + params.Url(), reqBody)

		if err != nil {

			sec <- err

			return
		}

		req.Header.Add("User-Agent", "bbacb081bc184538b1c8aa360036bd82")
		req.Header.Add("accept", "application/json")

		if params.RequiresAuth() {
			req.Header.Add("Authorization", "Bearer "+params.GetToken())
		}

		//if resp, err := http.Get(c.address + params.Url()); err == nil {

		if resp, err := client.Do(req); err == nil {

			defer resp.Body.Close()

			if resp.StatusCode != 200 {

				sec <- errors.New(fmt.Sprintf("RPC error, status %s", resp.Status))

				return

			}

			jd := json.NewDecoder(resp.Body)

			if errjd := jd.Decode(&result); errjd != nil {

				sec <- errors.New(fmt.Sprintf("Invalid json response: %s", errjd.Error()))

				return

			}

			//Cache result
			if exp, err := time.ParseInLocation(time.RFC1123, resp.Header.Get("expires"), time.Local); err == nil {

				if exp.After(time.Now()) {

					if b, errj := json.Marshal(result); errj == nil {

						registry.Registry.Session.Cache(c.address + params.Url(), b, time.Until(exp))

					}

				}

			}

			sec <- nil

		} else {

			sec <- err
		}

	}()

	sendErr = <-sec

	return sendErr

}
