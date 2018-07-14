// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/fstab/grok_exporter/tailer"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/oklog/run"
)

func main() {
	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen at %s: %s", addr, err)
		return
	}
	log.Printf("listening on %q", addr)

	var group run.Group
	{
		router := setUpRouter()
		group.Add(
			func() error { return http.Serve(l, router) },
			func(err error) { logErr(l.Close(), log.Print) })
	}
	{
		sigs := make(chan os.Signal)
		signal.Notify(sigs, os.Interrupt)
		group.Add(
			func() error {
				<-sigs
				log.Print("process interrupted")
				return nil
			},
			func(err error) {})
	}
	logErr(group.Run(), log.Fatal)
}

func setUpRouter() *mux.Router {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 1024,
		CheckOrigin:     func(_ *http.Request) bool { return true },
	}

	h := NewEventStreamingHandler(upgrader)

	r := mux.NewRouter()
	r.Path("/events").Handler(h)
	return r
}

type EventStreamingHandler struct {
	upgrader websocket.Upgrader
}

func NewEventStreamingHandler(upgrader websocket.Upgrader) *EventStreamingHandler {
	return &EventStreamingHandler{upgrader}
}

func (h *EventStreamingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	const (
		readAll       = true
		failOnMissing = false
	)
	tail := tailer.RunFseventFileTailer("./log", readAll, failOnMissing, nil)
	defer tail.Close()

	h.streamEvents(conn, tail)
}

func (h *EventStreamingHandler) streamEvents(conn *websocket.Conn, tail tailer.Tailer) {
	for {
		select {
		case err := <-tail.Errors():
			conn.WriteJSON(err.Error())
		case evt := <-tail.Lines():
			conn.WriteJSON(evt)
		}
	}
}

func logErr(err error, logF func(...interface{})) {
	if err == nil {
		return
	}
	logF(err)
}
