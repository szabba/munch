// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package handlers

import (
	"io"

	"github.com/szabba/munch"
)

func Discard() OnMessager {
	return discard{}
}

type discard struct{}

func (_ discard) OnMessage(_ munch.ClientID, _ io.Reader) {}
