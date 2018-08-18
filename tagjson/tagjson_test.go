// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tagjson_test

import (
	"bytes"
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/munch/tagjson"
)

func TestTagWithTypeOK(t *testing.T) {
	// given
	v := ""
	outWant := []byte(`{"string":""}`)

	// when
	out, err := tagjson.TagWithType(v)

	// then
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
	assert.That(bytes.Equal(out, outWant), t.Errorf, "got encoding %q, want %q", out, outWant)
}

func TestUntagFailsOnEmptyInput(t *testing.T) {
	// given
	in := []byte(``)

	// when
	tag, nested, err := tagjson.Untag(in)

	// then
	assert.That(err != nil, t.Fatalf, "got no error, wanted one")
	assert.That(tag == "", t.Errorf, "got tag %q, want %q", tag, "")
	assert.That(string(nested) == "", t.Errorf, "got nested message %q, want %q", nested, "")
}

func TestUntagFailsOnUntaggedInput(t *testing.T) {
	// given
	in := []byte(`{}`)

	// when
	tag, nested, err := tagjson.Untag(in)

	// then
	assert.That(err != nil, t.Fatalf, "got no error, wanted one")
	assert.That(tag == "", t.Errorf, "got tag %q, want %q", tag, "")
	assert.That(string(nested) == "", t.Errorf, "got nested message %q, want %q", nested, "")
}

func TestUntagFailsOnMultiKeyInput(t *testing.T) {
	// given
	in := []byte(`{ "x": 2, "y": null }`)

	// when
	tag, nested, err := tagjson.Untag(in)

	// then
	assert.That(err != nil, t.Fatalf, "got no error, wanted one")
	assert.That(tag == "", t.Errorf, "got tag %q, want %q", tag, "")
	assert.That(string(nested) == "", t.Errorf, "got nested message %q, want %q", nested, "")
}

func TestUntagSuccedsForSingleKeyInput(t *testing.T) {
	// given
	in := []byte(`{ "x": 2 }`)

	// when
	tag, nested, err := tagjson.Untag(in)

	// then
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
	assert.That(tag == "x", t.Errorf, "got tag %q, want %q", tag, "x")
	assert.That(string(nested) == "2", t.Errorf, "got nested message %q, want %q", nested, "2")
}
