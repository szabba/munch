// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers_test

import (
	"bytes"
	"fmt"
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

type SprintFormatter struct{}

var _ handlers.MessageFormatter = SprintFormatter{}

func (fmtr SprintFormatter) FormatMessage(w io.Writer, msg interface{}) error {
	_, err := fmt.Fprint(w, msg)
	return err
}

type CaptureFormatter struct {
	wasCalled bool
	nested    handlers.MessageFormatter
}

var _ handlers.MessageFormatter = new(CaptureFormatter)

func NewCaptureFormatter(nested handlers.MessageFormatter) *CaptureFormatter {
	return &CaptureFormatter{nested: nested}
}

func (fmtr *CaptureFormatter) FormatMessage(w io.Writer, msg interface{}) error {
	fmtr.wasCalled = true
	return fmtr.nested.FormatMessage(w, msg)
}

func (fmtr *CaptureFormatter) WasCalled() bool { return fmtr.wasCalled }
