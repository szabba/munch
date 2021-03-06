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

type BroadcastService interface {
	Broadcast(msg interface{})
}

type TailService struct {
	lock sync.Mutex
	once sync.Once

	tail    tailer.Tailer
	cast    BroadcastService
	clients map[munch.ClientID]chan<- interface{}
}

func NewTailService(tail tailer.Tailer, cast BroadcastService) *TailService {
	return &TailService{
		tail:    tail,
		cast:    cast,
		clients: make(map[munch.ClientID]chan<- interface{}),
	}
}

func (t *TailService) Stop() {
	t.tail.Close()
}

func (t *TailService) Run() error {
	for {
		select {
		case err := <-t.tail.Errors():
			// TODO: specify an encoding
			t.cast.Broadcast(err)
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
	t.cast.Broadcast(munch.Event{
		Source:  "tail",
		At:      time.Now(),
		Message: line,
	})
}
