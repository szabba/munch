// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package notification

import (
	"log"
	"sync"

	"github.com/szabba/munch"
)

type Service struct {
	lock    sync.Mutex
	clients map[munch.ClientID]chan<- interface{}
}

func NewService() *Service {
	return &Service{
		clients: make(map[munch.ClientID]chan<- interface{}),
	}
}

func (srv *Service) Subscribe(id munch.ClientID, sink chan<- interface{}) {
	srv.lock.Lock()
	if sink != nil {
		srv.clients[id] = sink
	}
	srv.lock.Unlock()
}

func (srv *Service) Send(id munch.ClientID, msg interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	sink := srv.clients[id]
	if sink == nil {
		log.Printf("got message for unubscribed client %s: %#v", id, msg)
		return
	}

	sink <- msg
}

func (srv *Service) Unsubscribe(id munch.ClientID) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	srv.unsubscribe(id)
}

func (srv *Service) Broadcast(v interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, sink := range srv.clients {
		sink <- v
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
	sink := srv.clients[id]
	if sink != nil {
		close(sink)
		delete(srv.clients, id)
	}
}
