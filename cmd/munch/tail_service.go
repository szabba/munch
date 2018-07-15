// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"sync"
	"time"

	"github.com/fstab/grok_exporter/tailer"

	"github.com/szabba/munch"
)

type TailService struct {
	lock sync.Mutex
	once sync.Once

	tail    tailer.Tailer
	clients map[munch.ClientID]chan<- interface{}
}

func NewTailService(tail tailer.Tailer) *TailService {
	return &TailService{
		tail:    tail,
		clients: make(map[munch.ClientID]chan<- interface{}),
	}
}

func (t *TailService) Stop() {
	t.tail.Close()

	t.lock.Lock()
	defer t.lock.Unlock()
	for id := range t.clients {
		t.unsubscribe(id)
	}
}

func (t *TailService) Run() error {
	for {
		select {
		case err := <-t.tail.Errors():
			// TODO: broadcast the error to the UI too
			return err
		case line, ok := <-t.tail.Lines():
			if !ok {
				log.Print("closed tailer")
				return nil
			}
			t.sendLine(line)
		}
	}
}

func (t *TailService) sendLine(line string) {
	t.broadcast(munch.Event{
		Source:  "tail",
		At:      time.Now(),
		Message: line,
	})
}

func (t *TailService) broadcast(m interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Printf("broadcasting %v", m)

	for _, sink := range t.clients {
		if sink != nil {
			sink <- m
		}
	}
}

func (t *TailService) Subscribe(id munch.ClientID, sink chan<- interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()

	_, registered := t.clients[id]
	if registered {
		log.Fatalf("duplicate registration of client %v", id)
		return
	}

	log.Printf("registering client %v", id)
	t.clients[id] = sink
}

func (t *TailService) Unsubscribe(id munch.ClientID) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.unsubscribe(id)
}

func (t *TailService) unsubscribe(id munch.ClientID) {
	sink, registered := t.clients[id]
	if !registered {
		return
	}
	log.Printf("unregistering client %v", id)
	close(sink)
	delete(t.clients, id)
}
