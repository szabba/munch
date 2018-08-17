// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sources

import (
	"io"
	"sync"

	"github.com/oklog/run"
)

type Source struct {
	once      sync.Once
	input     io.ReadCloser
	midWriter *io.PipeWriter
	midReader *io.PipeReader
	parser    io.WriteCloser
}

func New(input io.ReadCloser, parser io.WriteCloser) *Source {
	midReader, midWriter := io.Pipe()
	return &Source{
		input:     input,
		midWriter: midWriter,
		midReader: midReader,
		parser:    parser,
	}
}

func (src *Source) Process() error {
	g := new(run.Group)
	g.Add(
		func() error { return src.copy(src.midWriter, src.input) },
		func(err error) { src.midWriter.CloseWithError(err) })
	g.Add(
		func() error { return src.copy(src.parser, src.midReader) },
		func(err error) { src.midReader.CloseWithError(err) })
	return g.Run()
}

func (src *Source) copy(w io.Writer, r io.Reader) error {
	_, err := io.Copy(w, r)
	return err
}

func (src *Source) Stop() {
	src.once.Do(src.close)
}

func (src *Source) close() {
	src.input.Close()
	src.parser.Close()
}
