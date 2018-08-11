// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/szabba/assert"

	"github.com/szabba/munch"
	"github.com/szabba/munch/handlers"
)

func TestSocketSubscribesAndUsubscribesClient(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, _, subsMock := startMockServer(ctrl)
	defer srv.Close()

	clientID := munch.ClientIDOf(0)

	gomock.InOrder(
		subsMock.EXPECT().Subscribe(clientID, gomock.Any()),
		subsMock.EXPECT().Unsubscribe(clientID))

	// when
	conn, err := connect(srv)
	assumeNoError(t, err)
	defer conn.Close()
}

func TestSocketPassesClientMessageToHandler(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, onMsgMock, subsMock := startMockServer(ctrl)
	defer srv.Close()

	clientID := munch.ClientIDOf(0)

	send := make(chan string)
	msgGot := ""
	gomock.InOrder(
		subsMock.EXPECT().Subscribe(clientID, gomock.Any()),
		onMsgMock.EXPECT().OnMessage(clientID, gomock.Any()).
			Do(func(_ munch.ClientID, r io.Reader) {
				bs, _ := ioutil.ReadAll(r)
				send <- string(bs)
			}),
		subsMock.EXPECT().Unsubscribe(clientID))

	conn, err := connect(srv)
	assumeNoError(t, err)
	defer conn.Close()

	// when
	err = conn.WriteMessage(websocket.TextMessage, []byte("msg"))
	assertNoError(t, err)
	msgGot = <-send

	// then
	assert.That(msgGot == "msg", t.Errorf, "handler got message %q, want %q", msgGot, "msg")
}

func TestSocketLetsTheSubscriptionServiceSendMessagesToTheClient(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, _, subsMock := startMockServer(ctrl)
	defer srv.Close()

	clientID := munch.ClientIDOf(0)

	var (
		send func(interface{})
		sem  = make(chan func(interface{}))
	)

	gomock.InOrder(
		subsMock.EXPECT().Subscribe(clientID, gomock.Any()).
			Do(func(_ munch.ClientID, f func(interface{})) { sem <- f }),
		subsMock.EXPECT().Unsubscribe(clientID))

	conn, err := connect(srv)
	assumeNoError(t, err)
	defer conn.Close()

	send = <-sem
	// when
	send("msg")
	_, msgGot, err := conn.ReadMessage()
	assertNoError(t, err)

	// then
	assert.That(
		string(msgGot) == "msg",
		t.Errorf, "got message %q, want %q", msgGot, "msg")
}

func startMockServer(ctrl *gomock.Controller) (*httptest.Server, *MockOnMessager, *MockSubscriptionService) {
	onMsgMock := NewMockOnMessager(ctrl)
	subsMock := NewMockSubscriptionService(ctrl)
	srv := startServer(onMsgMock, subsMock)
	return srv, onMsgMock, subsMock
}

func startServer(onMsg handlers.OnMessager, subs handlers.SubscriptionService) *httptest.Server {
	up := websocket.Upgrader{}
	idGen := new(munch.ClientIDGenerator)
	h := handlers.NewSocket(up, idGen, onMsg, subs)
	return httptest.NewServer(h)
}

func connect(srv *httptest.Server) (*websocket.Conn, error) {
	wsURL, err := url.Parse(srv.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %s", srv.URL)
	}
	wsURL.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), make(http.Header))
	return conn, err
}

func assumeNoError(t *testing.T, err error) {
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
}

func assertNoError(t *testing.T, err error) {
	assert.That(err == nil, t.Errorf, "unexpected error: %s", err)
}
