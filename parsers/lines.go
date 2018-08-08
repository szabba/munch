// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"bytes"
	"time"

	"github.com/szabba/munch"
)

type Lines struct {
	clock func() time.Time
	cons  EventConsumer
	buf   bytes.Buffer
}

func NewLines(clock func() time.Time, cons EventConsumer) *Lines {
	return &Lines{clock: clock, cons: cons}
}

func (l *Lines) Write(p []byte) (n int, err error) {
	var line []byte
	for len(p) > 0 && err == nil {
		line, p, n = l.writeLine(p, n)
		if line == nil {
			continue
		}
		err := l.submitEvent(string(line))
		l.buf.Reset()
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (l *Lines) Close() error {
	line := l.buf.String()
	l.buf.Reset()
	return l.submitEvent(line)
}

func (l *Lines) writeLine(p []byte, n int) (line, left []byte, n2 int) {
	ix := bytes.IndexByte(p, '\n')
	if ix == -1 {
		l.buf.Write(p)
		return nil, nil, n + len(p)
	}
	prefix, suffix := p[:ix], p[ix+1:]
	l.buf.Write(prefix)
	return l.buf.Bytes(), suffix, n + len(prefix) + 1
}

func (l *Lines) submitEvent(msg string) error {
	now := l.clock()
	evt := munch.Event{At: now, Message: msg}
	return l.cons.On(evt)
}
