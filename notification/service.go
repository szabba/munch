// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package notification

import (
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

func (srv *Service) Unsubscribe(id munch.ClientID) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	sink := srv.clients[id]
	if sink != nil {
		delete(srv.clients, id)
	}
}

func (srv *Service) Broadcast(v interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, sink := range srv.clients {
		sink <- v
	}
}
