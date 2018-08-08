// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"time"

	"github.com/szabba/munch"
)

func stepClock(start time.Time, dt time.Duration) func() time.Time {
	next := start
	return func() time.Time {
		out := next
		next = out.Add(dt)
		return out
	}
}

type SliceConsumer struct {
	err  error
	evts []munch.Event
}

func (sc *SliceConsumer) On(evt munch.Event) error {
	if sc.err != nil {
		return sc.err
	}
	sc.evts = append(sc.evts, evt)
	return nil
}

func (sc *SliceConsumer) SetError(err error) { sc.err = err }

func (sc *SliceConsumer) Len() int { return len(sc.evts) }

func (sc *SliceConsumer) Event(i int) munch.Event { return sc.evts[i] }
