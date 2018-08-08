// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"fmt"
	"io"
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

type MessageHandler interface {
	OnMessage(id munch.ClientID, r io.Reader)
}

type SocketHandler struct {
	upgrader   websocket.Upgrader
	ids        ClientIDFactory
	msgHandler MessageHandler
	subs       SubscriptionService
}

func NewSocketHandler(upgrader websocket.Upgrader, ids ClientIDFactory, msgH MessageHandler, subs SubscriptionService) *SocketHandler {
	return &SocketHandler{
		upgrader:   upgrader,
		ids:        ids,
		msgHandler: noopIfNil(msgH),
		subs:       subs,
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
		_, r, err := conn.NextReader()
		h.msgHandler.OnMessage(id, r)
		if err != nil {
			log.Printf("client %v read error: %s", id, err)
			h.closeConn(conn, id)
			return err
		}
	}
}

func (h *SocketHandler) writeLoop(conn *websocket.Conn, id munch.ClientID, writes <-chan interface{}) error {
	for msg := range writes {
		err := h.writeMsg(conn, msg)
		if err != nil {
			log.Printf("client %v write error: %s", id, err)
			return err
		}
	}
	return nil
}

func (h *SocketHandler) writeMsg(conn *websocket.Conn, msg interface{}) error {
	wrapped := make(map[string]interface{})
	wrapped[fmt.Sprintf("%T", msg)] = msg
	return conn.WriteJSON(wrapped)
}

func (h *SocketHandler) closeConn(conn *websocket.Conn, id munch.ClientID) {
	err := conn.Close()
	if err != nil {
		log.Printf("client %v close error: %s", id, err)
	}
}

func noopIfNil(mh MessageHandler) MessageHandler {
	if mh == nil {
		return noopMsgHandler{}
	}
	return mh
}

type noopMsgHandler struct{}

func (_ noopMsgHandler) OnMessage(_ munch.ClientID, _ io.Reader) {}
