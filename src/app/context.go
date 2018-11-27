package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/flosch/pongo2"
	"net/http"
	"encoding/json"
	"sync"
	"time"
	"fmt"
	"runtime/debug"
)

type Context struct {
	http.ResponseWriter
	mutex    *sync.Mutex
	values   map[string]interface{}
	Request  *http.Request
	Params   httprouter.Params
	handlers []Handler
	index    int
	stop     bool
}

func (c *Context) GetUrlParam(key string) string {

	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {

	return c.Params.ByName(key)
}

func (c *Context) Set(key string, value interface{}) {

	c.mutex.Lock()

	c.values[key] = value

	c.mutex.Unlock()
}

func (c *Context) Get(key string) interface{} {

	defer c.mutex.Unlock()

	c.mutex.Lock()

	if v, found := c.values[key]; found {

		return v
	}

	return nil
}

func (c *Context) run() {


	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic:", r)
			debug.PrintStack()
			c.WriteHeader(http.StatusInternalServerError)
			c.Write([]byte(fmt.Sprintf("%v", r)))
		}
	}()

	go func() {

		select {

		case <-c.Request.Context().Done():

			c.stop = true

			return

		}

	}()

	for c.index < len(c.handlers) {

		c.handlers[c.index](c)

		c.index++

		if c.stop {

			return
		}
	}

}

func (c *Context) Stop() {

	c.stop = true
}

func (c *Context) IsPost() bool {

	if c.Request.Method == "POST" {

		return true

	}

	return false
}

func (c *Context) RenderJSON(value interface {}) error {

	select {

	case <-c.Request.Context().Done():

		c.stop = true

		return c.Request.Context().Err()

	default:

		c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		return json.NewEncoder(c.ResponseWriter).Encode(value)
	}

}

func (c *Context) RenderHTML(tmplName string, context map[string]interface {}) {

	if tn, ok := Templates[tmplName]; ok {

		select {

		case <-c.Request.Context().Done():

			c.stop = true

			return

		default:

			c.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
			c.ResponseWriter.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
			c.ResponseWriter.Header().Set("Pragma", "no-cache")
			c.ResponseWriter.Header().Set("Connection", "keep-alive")
			c.ResponseWriter.Header().Set("Expires", "0")

			context["currentUser"] = c.Get("user")
			context["currentTime"] = time.Now()
			context["lang"] = c.Get("lang")
			context["userSettings"] = c.Get("userSettings")
			context["pinnedBPOCount"] = c.Get("pinnedBPOCount")

			if err := tn.ExecuteWriter(pongo2.Context(context), c.ResponseWriter); err != nil {

				//c.RenderError(err.Error())

				c.RenderHTML("/errors/500.html", map[string]interface {} {
					"Error" : err.Error(),
				})

				c.Stop()

				return
			}

		}

	} else {

		c.RenderHTML("/errors/404.html", map[string]interface {} {})

		c.Stop()

		return
	}

}

func (c *Context) RenderError(e string) error {

	return json.NewEncoder(c.ResponseWriter).Encode(map[string]string{"Error" : e})

}

func (c *Context) Redirect(url string) {

	c.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header().Set("Expires", "0")

	http.Redirect(c.ResponseWriter, c.Request, url, 302)

	c.Stop()

	return
}