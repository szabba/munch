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

	sender := NewTestSender(t)

	service.Subscribe(ClientID, sender.Send)
	defer service.Unsubscribe(ClientID)
	msg := "msg"

	// when
	service.Broadcast(msg)

	// then
	sender.AssertGotString(msg)
}

func TestServiceDoesNotSendBroadcastToUnsubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	sender := NewTestSender(t)

	service.Subscribe(ClientID, sender.Send)
	service.Unsubscribe(ClientID)

	msg := "msg"

	// when
	service.Broadcast(msg)

	// then
	sender.AssertGotNothing()
}

func TestServiceSendsTargetedMessageToTheAddressee(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sender := NewTestSender(t)
	msg := "msg"

	service.Subscribe(id, sender.Send)
	defer service.Unsubscribe(id)

	// when
	service.Send(id, msg)

	// then
	sender.AssertGotString(msg)
}

func TestServiceDoesNotSendTargtedMessageToADifferentSubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	dstID, nonDstID := idGenerator.NextID(), idGenerator.NextID()
	dst, nonDst := NewTestSender(t), NewTestSender(t)

	service.Subscribe(dstID, dst.Send)
	defer service.Unsubscribe(dstID)
	service.Subscribe(nonDstID, nonDst.Send)
	defer service.Unsubscribe(nonDstID)

	msg := "msg"

	// when
	service.Send(dstID, msg)

	// then
	nonDst.AssertGotNothing()
	dst.AssertGotString(msg)
}

func TestServiceDoesNotSendTargetedMessageToAnUnsubscribedClient(t *testing.T) {
	// given
	service := notification.NewService()
	defer service.Close()

	idGenerator := new(munch.ClientIDGenerator)
	id := idGenerator.NextID()
	sender := NewTestSender(t)
	msg := "msg"

	service.Subscribe(id, sender.Send)
	service.Unsubscribe(id)

	// when
	service.Send(id, msg)

	// then
	sender.AssertGotNothing()
}

type TestSender struct {
	t       *testing.T
	wasSent bool
	msgGot  interface{}
}

func NewTestSender(t *testing.T) *TestSender {
	return &TestSender{t: t}
}

func (s *TestSender) Send(msg interface{}) {
	s.wasSent = true
	s.msgGot = msg
}

func (s *TestSender) AssertGotString(msgWant string) {
	s.t.Helper()

	assert.That(s.wasSent, s.t.Errorf, "no message was sent")
	if !s.wasSent {
		return
	}

	msgGot, isString := s.msgGot.(string)
	assert.That(isString, s.t.Errorf, "got message of type %T, want %T", s.msgGot, Message)
	if !isString {
		return
	}

	assert.That(msgGot == Message, s.t.Errorf, "got message %q, want %q", msgGot, Message)
}

func (s *TestSender) AssertGotNothing() {
	s.t.Helper()

	assert.That(!s.wasSent, s.t.Errorf, "unexpected message was sent: %#v", s.msgGot)
}
