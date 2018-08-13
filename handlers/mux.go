// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"encoding/json"
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

func (mux *Mux) OnMessage(id munch.ClientID, msg json.RawMessage) {

	tag, inner, ok := mux.decodeInnerMessage(id, msg)
	if !ok {
		return
	}

	h := mux.tagHandlers[tag]
	if h == nil {
		log.Printf("client %s sent message with unexpected tag: %s", id, tag)
		return
	}

	h.OnMessage(id, inner)
}

func (mux *Mux) decodeInnerMessage(id munch.ClientID, tagged json.RawMessage) (string, json.RawMessage, bool) {
	msg := make(map[string]json.RawMessage)
	err := json.Unmarshal(tagged, &msg)
	if err != nil {
		log.Printf("client %s sent invalid message: %s: %s", id, err, tagged)
		return "", nil, false
	}

	tag, inner, ok := mux.splitInnerMessage(msg)
	if !ok {
		log.Printf("client %s sent invalid message: %s", id, tagged)
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
