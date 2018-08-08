// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"net"
	"net/http"

	"github.com/szabba/munch/notification"

	"github.com/fstab/grok_exporter/tailer"
	"github.com/gorilla/websocket"
	"github.com/oklog/run"

	"github.com/szabba/munch"
	"github.com/szabba/munch/handlers"
)

func main() {
	addr := ":8080"

	interruptHandler := NewInterruptHandler()

	tail := NewFSTailer()
	defer tail.Close()

	notifSvc := notification.NewService()
	defer notifSvc.Close()

	tailService := NewTailService(tail, notifSvc)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(_ *http.Request) bool { return true },
	}

	clientIDGen := new(munch.ClientIDGenerator)

	sockHandler := handlers.NewSocketHandler(upgrader, clientIDGen, nil, notifSvc)

	l, err := net.Listen("tcp", addr)
	logErr(err, log.Fatal)
	log.Printf("listening on %q", addr)

	var group run.Group

	group.Add(interruptHandler.Run, func(_ error) { interruptHandler.Stop() })
	group.Add(tailService.Run, func(_ error) { tailService.Stop() })
	group.Add(
		func() error { return http.Serve(l, sockHandler) },
		func(_ error) { l.Close() },
	)

	logErr(group.Run(), log.Fatal)
}

func NewFSTailer() tailer.Tailer {
	const (
		readAll       = true
		failOnMissing = false
	)
	return tailer.RunFseventFileTailer("./log", readAll, failOnMissing, nil)
}

func logErr(err error, logF func(...interface{})) {
	if err == nil {
		return
	}
	logF(err)
}
