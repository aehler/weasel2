package app

import (
	"app/registry"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"runtime/debug"
	"log"
	"github.com/flosch/pongo2"
	"fmt"
	"golang.org/x/net/websocket"
	"app/bindata/assets"
	"gopkg.in/yaml.v2"
)

type App struct {
	Router *httprouter.Router
	handlers []Handler
}

type Handler404 struct {}

func New(configData []byte, config string, withBin bool) *App {

	a := App{
		Router: httprouter.New(),
	}

	rr := registry.Cfg{}

	if err := yaml.Unmarshal(configData, &rr); err != nil {

		log.Fatal(err.Error())
	}

	pathes := rr.Path

	log.Println("Using pathes", pathes)

	if withBin {

		InitBinaryTemplates(pathes.Templates)

		fmt.Printf("Serve static from bindata on /%s/*filepath\n", pathes.HTTPStatic)

		a.Router.ServeFiles(fmt.Sprintf("/%s/*filepath", pathes.HTTPStatic), assets.AssetFS())

		a.Router.ServeFiles(fmt.Sprintf("/%s/*filepath", pathes.ImageStatic), http.Dir(pathes.Static))

	} else {

		InitTemplates(pathes.Templates)

		fmt.Printf("Serve static on /%s/*filepath\n", pathes.HTTPStatic)

		a.Router.ServeFiles(fmt.Sprintf("/%s/*filepath", pathes.HTTPStatic), http.Dir(pathes.Static))

	}

	a.Router.NotFound = Handler404{}

	a.Router.PanicHandler = func(rw http.ResponseWriter, _ *http.Request, err interface{}) {

		rw.Header().Set("Content-Type", "text/html")

		rw.WriteHeader(http.StatusInternalServerError)

		Templates["/errors/500.html"].ExecuteWriter(pongo2.Context{"Error" : "DON'T PANIC", "Message" : fmt.Sprintf("%v\n", err)+string(debug.Stack())}, rw)

		log.Printf("PANIC: %v\n %s\n", err, debug.Stack())
	}

	registry.Init(rr, config)

	a.Get("/metrics/", metricsHandler)

	return &a
}

func (e Handler404) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {

	Templates["/errors/404.html"].ExecuteWriter(pongo2.Context{}, rw)

}

func (a *App) Get(route string, handlers ...Handler) {

	handler := handler(append(a.handlers, handlers...))

	a.Router.GET(route, handler)
	a.Router.HEAD(route, handler)

}

func (a *App) Post(pattern string, handlers ...Handler) {

	a.Router.POST(pattern, handler(append(a.handlers, handlers...)))
}

func (a *App) GetPost(pattern string, handlers ...Handler) {

	handler := handler(append(a.handlers, handlers...))

	a.Router.GET(pattern, handler)
	a.Router.POST(pattern, handler)
}

func (a *App) Handler(h Handler) {

	a.handlers = append(a.handlers, h)
}

func (a *App) WebSocket(pattern string, handler websocket.Handler) {

	http.Handle(pattern, handler)

}

func Redirect(url string, c *Context, code int) {

	c.Stop()

	http.Redirect(c.ResponseWriter, c.Request, url, code)

}
