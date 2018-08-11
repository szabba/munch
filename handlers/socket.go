// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/szabba/munch"
)

type ClientIDFactory interface {
	NextID() munch.ClientID
}

type SubscriptionService interface {
	Subscribe(munch.ClientID, func(interface{}))
	Unsubscribe(munch.ClientID)
}

type OnMessager interface {
	OnMessage(id munch.ClientID, r io.Reader)
}

type Socket struct {
	upgrader   websocket.Upgrader
	ids        ClientIDFactory
	msgHandler OnMessager
	subs       SubscriptionService
}

func NewSocket(
	upgrader websocket.Upgrader,
	ids ClientIDFactory,
	onMsg OnMessager,
	subs SubscriptionService,
) *Socket {

	return &Socket{
		upgrader:   upgrader,
		ids:        ids,
		msgHandler: onMsg,
		subs:       subs,
	}
}

func (h *Socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	id := h.ids.NextID()
	sndr := sender{id: id, conn: conn}
	h.subs.Subscribe(id, sndr.send)
	defer h.subs.Unsubscribe(id)

	h.readLoop(id, conn)
}

func (h *Socket) readLoop(id munch.ClientID, conn *websocket.Conn) {
	for {
		_, r, err := conn.NextReader()
		if err != nil {
			log.Printf("client %s read error: %s", id, err)
			return
		}
		h.msgHandler.OnMessage(id, r)
	}
}

type sender struct {
	lock sync.Mutex
	id   munch.ClientID
	conn *websocket.Conn
}

func (s sender) send(msg interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	payload := []byte(fmt.Sprint(msg))
	err := s.conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Printf("client %s write error: %s", s.id, err)
	}
}
