// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/szabba/munch"
	"github.com/szabba/munch/tagjson"
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
	tag, inner, err := tagjson.Untag(msg)
	if err != nil {
		log.Printf("client %s sent invalid message: %s", id, msg)
		return
	}

	h := mux.tagHandlers[tag]
	if h == nil {
		log.Printf("client %s sent message with unexpected tag: %s", id, tag)
		return
	}

	h.OnMessage(id, inner)
}
