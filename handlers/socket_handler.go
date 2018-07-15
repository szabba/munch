// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oklog/run"

	"github.com/szabba/munch"
)

type ClientIDFactory interface {
	NextID() munch.ClientID
}

type SubscriptionService interface {
	Subscribe(munch.ClientID, chan<- interface{})
	Unsubscribe(munch.ClientID)
}

type SocketHandler struct {
	upgrader websocket.Upgrader
	ids      ClientIDFactory
	subs     SubscriptionService
}

func NewSocketHandler(upgrader websocket.Upgrader, ids ClientIDFactory, subs SubscriptionService) *SocketHandler {
	return &SocketHandler{
		upgrader: upgrader,
		ids:      ids,
		subs:     subs,
	}
}

func (h *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	id := h.ids.NextID()
	writes := make(chan interface{})
	h.subs.Subscribe(id, writes)

	h.run(conn, id, writes)
}

func (h *SocketHandler) run(conn *websocket.Conn, id munch.ClientID, writes <-chan interface{}) {
	var g run.Group

	g.Add(
		func() error { return h.readLoop(conn, id) },
		func(_ error) { conn.Close() })

	g.Add(
		func() error { return h.writeLoop(conn, id, writes) },
		func(_ error) { h.subs.Unsubscribe(id) })

	g.Run()
}

func (h *SocketHandler) readLoop(conn *websocket.Conn, id munch.ClientID) error {
	for {
		_, _, err := conn.NextReader()
		if err != nil {
			log.Printf("client %v read error: %s", id, err)
			h.closeConn(conn, id)
			return err
		}
	}
}

func (h *SocketHandler) writeLoop(conn *websocket.Conn, id munch.ClientID, writes <-chan interface{}) error {
	for w := range writes {
		err := conn.WriteJSON(w)
		if err != nil {
			log.Printf("client %v write error: %s", id, err)
			return err
		}
	}
	return nil
}

func (h *SocketHandler) closeConn(conn *websocket.Conn, id munch.ClientID) {
	err := conn.Close()
	if err != nil {
		log.Printf("client %v close error: %s", id, err)
	}
}
