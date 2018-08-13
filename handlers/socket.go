// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"encoding/json"
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
	OnMessage(id munch.ClientID, msg json.RawMessage)
}

type MessageFormatter interface {
	FormatMessage(w io.Writer, msg interface{}) error
}

type Socket struct {
	upgrader websocket.Upgrader
	ids      ClientIDFactory
	onMsg    OnMessager
	fmtr     MessageFormatter
	subs     SubscriptionService
}

func NewSocket(
	upgrader websocket.Upgrader,
	ids ClientIDFactory,
	onMsg OnMessager,
	fmtr MessageFormatter,
	subs SubscriptionService,
) *Socket {

	return &Socket{upgrader, ids, onMsg, fmtr, subs}
}

func (h *Socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	id := h.ids.NextID()
	sndr := newSender(id, h.fmtr, conn)
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
		var msg json.RawMessage
		err = json.NewDecoder(r).Decode(&msg)
		if err != nil {
			log.Printf("client %s sent invalid message: %s", id, err)
			break
		}
		h.onMsg.OnMessage(id, msg)
	}
}

type sender struct {
	lock sync.Mutex
	id   munch.ClientID
	fmtr MessageFormatter
	conn *websocket.Conn
}

func newSender(id munch.ClientID, fmtr MessageFormatter, conn *websocket.Conn) *sender {
	return &sender{id: id, conn: conn, fmtr: fmtr}
}

func (s *sender) send(msg interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	w, err := s.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Printf("client %s write error: %s", s.id, err)
		return
	}
	defer w.Close()
	err = s.fmtr.FormatMessage(w, msg)
	if err != nil {
		log.Printf("client %s write error: %s", s.id, err)
	}
}
