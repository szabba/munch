// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sources_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/oklog/run"
	"github.com/szabba/assert"

	"github.com/szabba/munch/sources"
)

const SleepTime time.Duration = 100 * time.Millisecond

func TestSourceClosesInputWhenClosed(t *testing.T) {
	// given
	in, inWriter := io.Pipe()
	_, out := io.Pipe()

	src := sources.New(in, out)

	// when
	src.Stop()

	// then
	n, err := io.WriteString(inWriter, "123")
	assert.That(err != nil, t.Errorf, "got no error trying to write to closed source")
	assert.That(n == 0, t.Errorf, "wrote %d bytes, expected %d", n, 0)
}

func TestSourceClosesParserWhenClosed(t *testing.T) {
	// given
	in, _ := io.Pipe()
	outReader, out := io.Pipe()

	src := sources.New(in, out)

	// when
	src.Stop()

	// then
	allOut, _ := ioutil.ReadAll(outReader)
	assert.That(bytes.Equal(allOut, nil), t.Errorf, "read %q, expected %q", allOut, nil)
}

func TestSourceCopiesAllInputToParserWhenAllowedTo(t *testing.T) {
	in, inWriter := io.Pipe()
	outReader, out := io.Pipe()

	src := sources.New(in, out)

	buf := new(bytes.Buffer)

	g := new(run.Group)
	g.Add(
		src.Process,
		func(_ error) { time.Sleep(SleepTime); src.Stop() })
	g.Add(
		func() error { _, err := io.Copy(buf, outReader); return err },
		func(_ error) { time.Sleep(SleepTime); outReader.Close() })

	// when
	g.Add(
		func() error { _, err := io.WriteString(inWriter, "abcd"); return err },
		func(_ error) { inWriter.Close() })
	g.Run()

	// then
	assert.That(buf.String() == "abcd", t.Errorf, "got %q into parser, want %q", buf.String(), "abcd")
}

func TestSourceClosesUpWhenTheInputGetsExhausted(t *testing.T) {
	// given
	in, inWriter := io.Pipe()
	outReader, out := io.Pipe()

	src := sources.New(in, out)

	buf := new(bytes.Buffer)

	g := new(run.Group)
	g.Add(
		src.Process,
		func(_ error) { src.Stop() })
	g.Add(
		func() error { _, err := io.Copy(buf, outReader); return err },
		func(_ error) { outReader.Close() })

	// when
	g.Add(
		func() error { time.Sleep(SleepTime); return inWriter.Close() },
		func(_ error) {})
	g.Run()

	// then
	assert.That(buf.String() == "", t.Errorf, "got %q into parser, want %q", buf.String(), "")
}

func TestSourceClosesUpWhenTheParserRefusesFurtherInput(t *testing.T) {
	// given
	in, inWriter := io.Pipe()
	outReader, out := io.Pipe()

	src := sources.New(in, out)

	var n int

	g := new(run.Group)
	g.Add(
		src.Process,
		func(_ error) { src.Stop() })
	g.Add(
		func() error {
			var err error
			time.Sleep(SleepTime)
			_, err = io.WriteString(inWriter, "some input")
			return err
		},
		func(_ error) {})

	// when
	g.Add(
		func() error { return outReader.Close() },
		func(_ error) {})
	g.Run()

	// then
	assert.That(n == 0, t.Errorf, "wrote %d bytes, expected %d", n, 0)
}
