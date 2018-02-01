package index

import (
	"app"
	"app/registry"
	"lib/stats"
	"golang.org/x/net/websocket"
	"log"
	"encoding/json"
	"time"
	"strings"
	"lib/triggers"
)

func Index(c *app.Context) {

	host := strings.Split(c.Request.Host, ":")

	c.RenderHTML("/index.html", map[string]interface {} {
		"gs"  :   stats.GenStats,
		"host":   host[0],
		"wsport": registry.Registry.KVS.Get([]byte("ws_port")),
	})

}

func restartService(c *app.Context) {

	id := c.Params.ByName("service")

	err := triggers.RestartService(id)

	log.Println("Autorestart", id, ", error:", err.Error())

	c.RenderJSON(err.Error())

}

func CPU(c *app.Context) {

	v := registry.Registry.Redis.Get("CPU", 0, 1)

	c.RenderJSON(v)
}

func PIDS(c *app.Context) {

	c.RenderJSON(stats.Pids)

}

func Health(ws *websocket.Conn, sc chan bool) {

	id := ws.RemoteAddr().String() + "-" + ws.Request().RemoteAddr + "-" + ws.Request().UserAgent()

	var scc = make(map[string]chan bool)

	defer func() {

		for _, c := range scc {
			c <- true
		}

	}()

	for name, c := range registry.Registry.WsChan.Collection(id) {

		scc[name] = make(chan bool)

		go func(scc chan bool, name, id string, c chan interface{}) {

			log.Println("Added ws listenter, chanel ", name, c)

			for {

				select {

				case <-scc:

					log.Println("Closing channel listener", name)

					return

				case p := <-c:

					//log.Println("sending to websocket", name, id)

					if wsc, err := json.Marshal(map[string]interface{}{name: p}); err == nil {

						_, err2 := ws.Write(wsc)

						if err2 != nil {

							log.Println(err2.Error())

						}

					} else {

						log.Println("Marshal error", err.Error())

					}

					time.Sleep(time.Millisecond * 30)

					continue

				default:

					time.Sleep(time.Millisecond * 800)

					continue

				}

			}
		}(scc[name], name, id, c)

	}

	select {
	case <- sc:

		log.Println("Closing health, id", id)

		return
	}

}

//func Health2(ws *websocket.Conn) {
//
//	id := ws.RemoteAddr().String() + "-" + ws.Request().RemoteAddr + "-" + ws.Request().UserAgent()
//
//	for {
//
//		for name, c := range registry.Registry.WsChan.Collection(id) {
//
//			log.Println(name, id)
//
//			select {
//
//			case <- scc:
//
//				log.Println("Closing", ws)
//
//				return
//
//			case p := <- c:
//
//				log.Println("sending to websocket", name, id)
//
//				if wsc, err := json.Marshal(map[string]interface{}{name: p}); err == nil {
//
//					_, err2 := ws.Write(wsc)
//
//					if err2 != nil {
//
//						log.Println(err2.Error())
//
//					}
//
//				} else {
//
//					log.Println("Marshal error", err.Error())
//
//				}
//
//				continue
//
//			default:
//
//				//time.Sleep(time.Millisecond * 300)
//
//				continue
//
//			}
//
//		}
//
//	}
//
//	log.Println("Out of cycle")
//
//}