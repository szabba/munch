// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"
	"testing"
	"time"

	"github.com/szabba/assert"
	"github.com/szabba/munch/parsers"
)

func TestLinesDoNothingForEmptyInput(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	// when
	n, err := lines.Write(nil)

	// then
	assert.That(n == 0, t.Errorf, "got %d bytes read, want %d", n, 0)
	assert.That(err == nil, t.Errorf, "unexpected error %q", err)
	assert.That(cons.Len() == 0, t.Errorf, "got %d events submitted, want %d", cons.Len(), 0)
	for i := 0; i < cons.Len(); i++ {
		t.Logf("unexpected event %d: %v", i, cons.Event(i))
	}
}

func TestLinesAcceptsLinePrefixWithoutSubmittingAnEvent(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	// when
	n, err := lines.Write([]byte("abba"))

	// then
	assert.That(n == 4, t.Errorf, "got %d bytes read, want %d", n, 4)
	assert.That(err == nil, t.Errorf, "unexpected error %q", err)
	assert.That(cons.Len() == 0, t.Errorf, "got %d events submitted, want %d", cons.Len(), 0)
	for i := 0; i < cons.Len(); i++ {
		t.Logf("unexpected event %d: %v", i, cons.Event(i))
	}
}

func TestLinesSubmitsEventFromSingleWrite(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	// when
	n, err := lines.Write([]byte("abba\n"))

	// then
	assert.That(n == 5, t.Errorf, "got %d bytes read, want %d", n, 5)
	assert.That(err == nil, t.Errorf, "unexpected error %q", err)
	assert.That(cons.Len() == 1, t.Fatalf, "got %d events submitted, want %d", cons.Len(), 1)

	evt := cons.Event(0)
	assert.That(evt.At == time.Unix(0, 0), t.Errorf, "got event at %v, want %v", evt.At, time.Unix(0, 0))
	assert.That(evt.Message == "abba", t.Errorf, "got event message %q, want %q", evt.Message, "abba")
}

func TestLinesSubmitsEventMergedFromMultipleWrites(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	lines.Write([]byte("abba"))

	// when
	n, err := lines.Write([]byte("hanna\n"))

	// then
	assert.That(n == 6, t.Errorf, "got %d bytes read, want %d", n, 6)
	assert.That(err == nil, t.Errorf, "unexpected error %q", err)
	assert.That(cons.Len() == 1, t.Fatalf, "got %d events submitted, want %d", cons.Len(), 1)

	evt := cons.Event(0)
	assert.That(evt.At == time.Unix(0, 0), t.Errorf, "got event at %v, want %v", evt.At, time.Unix(0, 0))
	assert.That(evt.Message == "abbahanna", t.Errorf, "got event message %q, want %q", evt.Message, "abbahanna")
}

func TestLinesSubmitsMultipleEventsFromSingleWrite(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	// when
	n, err := lines.Write([]byte("abba\nhanna\n"))

	// then
	assert.That(n == 11, t.Errorf, "got %d bytes read, want %d", n, 11)
	assert.That(err == nil, t.Errorf, "unexpected error %q", err)
	assert.That(cons.Len() == 2, t.Fatalf, "got %d events submitted, want %d", cons.Len(), 1)

	first := cons.Event(0)
	assert.That(first.At == time.Unix(0, 0), t.Errorf, "got first event at %v, want %v", first.At, time.Unix(0, 0))
	assert.That(first.Message == "abba", t.Errorf, "got first event message %q, want %q", first.Message, "abba")

	second := cons.Event(1)
	assert.That(second.At == time.Unix(1, 0), t.Errorf, "got second event at %v, want %v", second.At, time.Unix(1, 0))
	assert.That(second.Message == "hanna", t.Errorf, "got second event message %q, want %q", second.Message, "hanna")
}

func TestLinesReportsTheConsumerError(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	errWant := errors.New("consumer error")
	cons.SetError(errWant)

	// when
	n, err := lines.Write([]byte("hanna\n"))

	// then
	assert.That(n == 6, t.Errorf, "got %d bytes read, want %d", n, 6)
	assert.That(err != nil, t.Fatalf, "no error reported")
	assert.That(err == errWant, t.Errorf, "got err %q, want %q", err, errWant)
}

func TestLinesReportsTheRetainedInputAsEventUponBeingClosed(t *testing.T) {
	// given
	clock := stepClock(time.Unix(0, 0), time.Second)
	var cons SliceConsumer
	lines := parsers.NewLines(clock, &cons)

	lines.Write([]byte("abba"))

	// when
	err := lines.Close()

	// then
	assert.That(err == nil, t.Errorf, "unexpected error: %s", err)
	assert.That(cons.Len() == 1, t.Fatalf, "got %d events submitted, want %d", cons.Len(), 1)

	evt := cons.Event(0)
	assert.That(evt.At.Equal(time.Unix(0, 0)), t.Errorf, "got event at %v, want %v", evt.At, time.Unix(0, 0))
	assert.That(evt.Message == "abba", t.Errorf, "got event message %q, want %q", evt.Message, "abba")
}
