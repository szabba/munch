// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/munch"

	"github.com/szabba/munch/handlers"
)

type MsgA struct{}
type MsgB struct{}

var ClientID = new(munch.ClientIDGenerator).NextID()

func TestMuxDelegatesToMatchingHandler(t *testing.T) {
	// given
	matching := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): matching,
		reflect.TypeOf(MsgB{}): handlers.Discard(),
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`{"handlers_test.MsgA": {}}`))

	// then
	assert.That(matching.WasCalled(), t.Fatalf, "matching handler was not called")
	assert.That(matching.ID() == ClientID, t.Errorf, "got client ID %s, want %s", matching.ID(), ClientID)
	assert.That(
		matching.MessageText() == `{}`,
		t.Errorf, "got inner message %q, want %q", matching.MessageText(), `{}`)
}

func TestMuxDoesNotCallNonMatchingHandler(t *testing.T) {
	// given
	matching := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): handlers.Discard(),
		reflect.TypeOf(MsgB{}): matching,
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`{"handlers_test.MsgA": {}}`))

	// then
	assert.That(!matching.WasCalled(), t.Fatalf, "non-matching handler was called")
}

func TestMuxDoesNotCallHandlerForInvalidMessage(t *testing.T) {
	// given
	capt := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): capt,
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`{{{`))

	// then
	assert.That(!capt.WasCalled(), t.Fatalf, "inner handler was unexpectedly called")
}

func TestMuxDoesNotCallHandlerForANonObjectMessage(t *testing.T) {
	// given
	capt := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): capt,
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`[1, 2, 3]`))

	// then
	assert.That(!capt.WasCalled(), t.Fatalf, "inner handler was unexpectedly called")
}

func TestMuxDoesNotCallHandlerForAMultipleFieldObjectMessage(t *testing.T) {
	// given
	capt := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): capt,
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`{"handlers_test.MsgA": {}, "handlers_test.MsgB": {}}`))

	// then
	assert.That(!capt.WasCalled(), t.Fatalf, "inner handler was unexpectedly called")
}

func TestMuxIgnoresMessageWithUnknownTag(t *testing.T) {
	// given
	capt := new(CaptureHandler)

	mux := handlers.NewMux(map[reflect.Type]handlers.OnMessager{
		reflect.TypeOf(MsgA{}): capt,
	})

	// when
	mux.OnMessage(ClientID, json.RawMessage(`{"handlers_test.MsgC": {}}`))

	// then
	assert.That(!capt.WasCalled(), t.Fatalf, "inner handler was unexpectedly called")
}
