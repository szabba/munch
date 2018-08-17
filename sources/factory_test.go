// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sources_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/szabba/ahm/assert"

	"github.com/szabba/munch/sources"
)

func TestFactoryFailsToCreateSourceIfAnInputCannotBeCreated(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputFactory := NewMockInputFactory(ctrl)
	parserFactory := NewMockParserFactory(ctrl)

	factory := sources.NewFactory(inputFactory, parserFactory)

	def := sources.Definition{
		InputDefinition: json.RawMessage(`{}`),
		ParserDefition:  json.RawMessage(`null`),
	}

	errWant := errors.New("cannot create input")

	inputFactory.EXPECT().NewInput(def.InputDefinition).Return(nil, errWant)

	// when
	src, err := factory.NewSource(def)

	// then
	assert.That(err == errWant, t.Errorf, "got error %q, want %q", err, errWant)
	assert.That(src == nil, t.Errorf, "got source %#v, want %#v", src, nil)
}

func TestFactoryFailsToCreateSourceIfAParserCannotBeCreated(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputFactory := NewMockInputFactory(ctrl)
	parserFactory := NewMockParserFactory(ctrl)

	factory := sources.NewFactory(inputFactory, parserFactory)

	def := sources.Definition{
		InputDefinition: json.RawMessage(`{}`),
		ParserDefition:  json.RawMessage(`null`),
	}

	errWant := errors.New("cannot create parser")

	input := NewMockReadCloser(ctrl)

	gomock.InOrder(
		inputFactory.EXPECT().NewInput(def.InputDefinition).Return(input, nil),
		parserFactory.EXPECT().NewParser(def.ParserDefition).Return(nil, errWant),
		input.EXPECT().Close())

	// when
	src, err := factory.NewSource(def)

	// then
	assert.That(err == errWant, t.Errorf, "got error %q, want %q", err, errWant)
	assert.That(src == nil, t.Errorf, "got source %#v, want %#v", src, nil)
}

func TestFactoryCreatesSourceIfBothInputAndParserCanBeCreated(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputFactory := NewMockInputFactory(ctrl)
	parserFactory := NewMockParserFactory(ctrl)

	factory := sources.NewFactory(inputFactory, parserFactory)

	def := sources.Definition{
		InputDefinition: json.RawMessage(`{}`),
		ParserDefition:  json.RawMessage(`null`),
	}

	input := NewMockReadCloser(ctrl)
	parser := NewMockWriteCloser(ctrl)

	gomock.InOrder(
		inputFactory.EXPECT().NewInput(def.InputDefinition).Return(input, nil),
		parserFactory.EXPECT().NewParser(def.ParserDefition).Return(parser, nil))

	// when
	src, err := factory.NewSource(def)

	// then
	assert.That(err == nil, t.Errorf, "unexpected error: %s", err)
	assert.That(src != nil, t.Errorf, "got source %#v, want %#v", src, nil)
}
