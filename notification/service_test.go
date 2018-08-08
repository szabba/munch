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

const Message = "msg"

var ClientID = new(munch.ClientIDGenerator).NextID()

func TestServiceSendsBroadcastToSubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	sink := make(chan interface{}, 1)

	service.Subscribe(ClientID, sink)
	defer service.Unsubscribe(ClientID)

	// when
	service.Broadcast(Message)

	// then
	assertGotStringMessage(t, sink, Message)
}

func TestServiceDoesNotSendBroadcastToUnsubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	sink := make(chan interface{}, 1)

	service.Subscribe(ClientID, sink)
	service.Unsubscribe(ClientID)

	// when
	service.Broadcast(Message)

	// then
	assertNoMessagesSent(t, sink)
}

func TestServiceSendsTargetedMessageToTheAddressee(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)
	msg := "msg"

	service.Subscribe(id, sink)
	defer service.Unsubscribe(id)

	// when
	service.Send(id, msg)

	// then
	assertGotStringMessage(t, sink, Message)
}

func TestServiceDoesNotSendTargtedMessageToADifferentSubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	dstID, nonDstID := idGenerator.NextID(), idGenerator.NextID()
	dst := make(chan interface{}, 1)
	nonDst := make(chan interface{}, 1)

	service.Subscribe(dstID, dst)
	defer service.Unsubscribe(dstID)
	service.Subscribe(nonDstID, nonDst)
	defer service.Unsubscribe(nonDstID)

	// when
	service.Send(dstID, Message)

	// then
	assertGotStringMessage(t, dst, Message)
	assertNoMessagesSent(t, nonDst)
}

func TestServiceDoesNotSendTargetedMessageToAnUnsubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)
	msg := "msg"

	service.Subscribe(id, sink)
	service.Unsubscribe(id)

	// when
	service.Send(id, msg)

	// then
	assertNoMessagesSent(t, sink)
}

func TestServiceClosesSinkWhenUnsubscribingItsClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)

	service.Subscribe(id, sink)

	// when
	service.Unsubscribe(id)

	// then
	assertChannelWasClosed(t, sink)
}

func TestServiceClosesSinkWhenStoping(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sink := make(chan interface{}, 1)

	service.Subscribe(id, sink)

	// when
	service.Close()

	// then
	assertChannelWasClosed(t, sink)
}

func assertGotStringMessage(t *testing.T, ch chan interface{}, wantMsg string) {
	t.Helper()

	var rawMsg interface{}
	select {
	case rawMsg = <-ch:
	default:
	}

	msgGot, isString := rawMsg.(string)
	assert.That(isString, t.Fatalf, "got message of type %T, want %T", rawMsg, Message)
	assert.That(msgGot == Message, t.Errorf, "got message %q, want %q", msgGot, Message)
}

func assertNoMessagesSent(t *testing.T, ch chan interface{}) {
	t.Helper()

	msgCount := len(ch)
	assert.That(msgCount == 0, t.Errorf, "got %d messages in channel, wanted %d", msgCount, 0)
	for i := 0; i < msgCount; i++ {
		msg := <-ch
		t.Logf("unexpected message: %#v", msg)
	}
}

func assertChannelWasClosed(t *testing.T, ch chan interface{}) {
	t.Helper()

	isOpen := false
	select {
	case _, isOpen = <-ch:
	default:
	}

	assert.That(!isOpen, t.Errorf, "the channel was not closed")
}
