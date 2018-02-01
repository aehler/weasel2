package index

import (
	"app/registry"
	"golang.org/x/net/websocket"
	"log"
)

func WsHandler(ws *websocket.Conn) {

	id := ws.RemoteAddr().String() + "-" + ws.Request().RemoteAddr + "-" + ws.Request().UserAgent()

	registry.Registry.WsChan.RegisterCollection(id)

	registry.Registry.WsChan.RegisterInCollection("pids", id)
	registry.Registry.WsChan.RegisterInCollection("cpuPercent", id)
	registry.Registry.WsChan.RegisterInCollection("ps", id)
	registry.Registry.WsChan.RegisterInCollection("openTCPSockets", id)
	registry.Registry.WsChan.RegisterInCollection("serviceFailure", id)
	registry.Registry.WsChan.RegisterInCollection("memStat", id)

	var scc = make(chan bool)

	defer func() {

		log.Println("Connection closed", id)

		registry.Registry.WsChan.UnsetCollection(id)

		ws.Close()

		scc <- true

	}()

	log.Println("New webSocket open", id)

	go Health(ws, scc)

	readPump(ws)

}

func readPump(ws *websocket.Conn) {

	defer ws.Close()

	msg := make([]byte, 512)
	_, err := ws.Read(msg)
	if err != nil {
		return
	}

	log.Println(msg)
}