// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/szabba/munch/handlers"
)

type TagFormatter struct{}

var _ handlers.MessageFormatter = TagFormatter{}

func (fmtr TagFormatter) FormatMessage(w io.Writer, msg interface{}) error {
	tagged := map[string]interface{}{
		fmt.Sprintf("%T", msg): msg,
	}
	enc := json.NewEncoder(w)
	return enc.Encode(tagged)
}
