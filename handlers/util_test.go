// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers_test

import (
	"bytes"
	"io"

	"github.com/szabba/munch"
	"github.com/szabba/munch/handlers"
)

type CaptureHandler struct {
	wasCalled bool
	id        munch.ClientID
	buf       bytes.Buffer
}

var _ handlers.OnMessager = new(CaptureHandler)

func (capt *CaptureHandler) OnMessage(id munch.ClientID, r io.Reader) {
	capt.wasCalled = true
	capt.id = id
	capt.buf.Reset()
	io.Copy(&capt.buf, r)
}

func (capt *CaptureHandler) WasCalled() bool     { return capt.wasCalled }
func (capt *CaptureHandler) ID() munch.ClientID  { return capt.id }
func (capt *CaptureHandler) MessageText() string { return capt.buf.String() }
