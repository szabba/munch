// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package notification_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/munch"
	"github.com/szabba/munch/notification"
)

func TestServiceSendsBroadcastToSubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)
	msg := "msg"

	service.Subscribe(id, sink)
	defer service.Unsubscribe(id)

	// when
	service.Broadcast(msg)
	rawMsgReceived := <-sink

	// then
	msgReceived, isString := rawMsgReceived.(string)
	assert.That(isString, t.Fatalf, "got message of type %T, want %T", rawMsgReceived, msg)
	assert.That(msgReceived == msg, t.Errorf, "got message %q, want %q", msgReceived, msg)
}

func TestServiceDoesNotSentBroadcastToUnsubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)
	msg := "msg"

	service.Subscribe(id, sink)
	service.Unsubscribe(id)

	// when
	service.Broadcast(msg)
	msgCount := len(sink)

	// then
	assert.That(msgCount == 0, t.Errorf, "got %d messages in channel, wanted %d", msgCount, 0)
	for i := 0; i < msgCount; i++ {
		msg := <-sink
		t.Logf("unexpected message: %#v", msg)
	}
}

func TestServiceClosesSinkWhenUnsubscribingItsClient(t *testing.T) {
	// given
	service := notification.NewService()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)

	service.Subscribe(id, sink)

	// when
	service.Unsubscribe(id)

	// then
	isOpen := false
	select {
	case _, isOpen = <-sink:
	default:
	}
	assert.That(!isOpen, t.Errorf, "the sink was not closed")
}

func TestServiceClosesSinkWhenStoping(t *testing.T) {
	// given
	service := notification.NewService()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)

	service.Subscribe(id, sink)

	// when
	service.Close()

	// then
	isOpen := false
	select {
	case _, isOpen = <-sink:
	default:
	}
	assert.That(!isOpen, t.Errorf, "the sink was not closed")
}
