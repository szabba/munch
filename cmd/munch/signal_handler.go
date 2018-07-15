// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"os"
	"os/signal"
)

type InterruptHandler struct {
	sigs chan os.Signal
}

func NewInterruptHandler() *InterruptHandler {
	return &InterruptHandler{
		sigs: make(chan os.Signal, 1),
	}
}

func (h *InterruptHandler) Run() error {
	signal.Notify(h.sigs, os.Interrupt)
	_, ok := <-h.sigs
	if ok {
		log.Printf("received interrupt signal")
	}
	return nil
}

func (h *InterruptHandler) Stop() {
	signal.Stop(h.sigs)
	close(h.sigs)
}
