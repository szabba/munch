// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sources

import (
	"encoding/json"
	"io"
)

type Factory struct {
	inputFactory  InputFactory
	parserFactory ParserFactory
}

type InputFactory interface {
	NewInput(json.RawMessage) (io.ReadCloser, error)
}

type ParserFactory interface {
	NewParser(json.RawMessage) (io.WriteCloser, error)
}

func NewFactory(inputFactory InputFactory, parserFactory ParserFactory) *Factory {
	return &Factory{
		inputFactory,
		parserFactory,
	}
}

func (fact *Factory) NewSource(def Definition) (*Source, error) {
	input, err := fact.inputFactory.NewInput(def.InputDefinition)
	if err != nil {
		return nil, err
	}
	parser, err := fact.parserFactory.NewParser(def.ParserDefition)
	if err != nil {
		input.Close()
		return nil, err
	}
	return New(input, parser), nil
}
