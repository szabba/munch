// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package notification

import (
	"log"
	"sync"

	"github.com/szabba/assert"

	"github.com/szabba/munch"
)

type Service struct {
	lock    sync.Mutex
	clients map[munch.ClientID]sender
}

type sender func(interface{})

func (s sender) send(msg interface{}) {
	s(msg)
}

func NewService() *Service {
	return &Service{
		clients: make(map[munch.ClientID]sender),
	}
}

func (srv *Service) Subscribe(id munch.ClientID, sndr func(interface{})) {
	srv.lock.Lock()
	assert.That(sndr != nil, log.Panicf, "client %s registration attempted with nil sender", id)
	srv.clients[id] = sndr
	srv.lock.Unlock()
}

func (srv *Service) Send(id munch.ClientID, msg interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	sender := srv.clients[id]
	if sender == nil {
		log.Printf("got message for unubscribed client %s: %#v", id, msg)
		return
	}
	sender.send(msg)
}

func (srv *Service) Unsubscribe(id munch.ClientID) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	srv.unsubscribe(id)
}

func (srv *Service) Broadcast(v interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, sender := range srv.clients {
		sender.send(v)
	}
}

func (srv *Service) Close() {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	ids := make([]munch.ClientID, 0, len(srv.clients))
	for id := range srv.clients {
		ids = append(ids, id)
	}
	for _, id := range ids {
		srv.unsubscribe(id)
	}
}

func (srv *Service) unsubscribe(id munch.ClientID) {
	delete(srv.clients, id)
}
