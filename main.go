// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"net/http"

	"github.com/fstab/grok_exporter/tailer"
	"github.com/gorilla/websocket"
)

func main() {
	h := &EventStreamingHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  0,
			WriteBufferSize: 1024,
			CheckOrigin:     func(_ *http.Request) bool { return true },
		},
	}
	log.Printf("listening on %q", ":8080")
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Print(err)
	}
}

type EventStreamingHandler struct {
	upgrader websocket.Upgrader
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

func logErr(err error) {
	if err == nil {
		return
	}
	log.Print(err)
}
