// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"reflect"

	"github.com/szabba/munch"
)

type Mux struct {
	tagHandlers map[string]OnMessager
}

func NewMux(hs map[reflect.Type]OnMessager) *Mux {
	tagHandlers := make(map[string]OnMessager)
	for typ, h := range hs {
		tagHandlers[typ.String()] = h
	}
	return &Mux{tagHandlers}
}

var _ OnMessager = new(Mux)

func (mux *Mux) OnMessage(id munch.ClientID, r io.Reader) {

	tag, inner, ok := mux.decodeInnerMessage(id, r)
	if !ok {
		return
	}

	h := mux.tagHandlers[tag]
	if h == nil {
		log.Printf("client %s sent message with unexpected tag: %s", id, tag)
		return
	}

	h.OnMessage(id, bytes.NewReader(inner))
}

func (mux *Mux) decodeInnerMessage(id munch.ClientID, r io.Reader) (string, json.RawMessage, bool) {
	msg := make(map[string]json.RawMessage)
	cp := new(bytes.Buffer)
	dec := json.NewDecoder(io.TeeReader(r, cp))

	err := dec.Decode(&msg)
	if err != nil {
		log.Printf("client %s sent invalid message: %s: %s", id, err, cp.Bytes())
		return "", nil, false
	}

	tag, inner, ok := mux.splitInnerMessage(msg)
	if !ok {
		log.Printf("client %s sent invalid message: %s", id, cp.Bytes())
		return "", nil, false
	}
	return tag, inner, true
}

func (mux *Mux) splitInnerMessage(msg map[string]json.RawMessage) (string, json.RawMessage, bool) {
	if len(msg) != 1 {
		return "", nil, false
	}

	for tag, in := range msg {
		return tag, in, true
	}

	// Unreachable.
	return "", nil, false
}
