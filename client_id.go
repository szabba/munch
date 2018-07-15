// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package munch

import (
	"fmt"
	"sync"
)

type ClientID struct{ id int }

func (id ClientID) String() string { return fmt.Sprint(id.id) }

type ClientIDGenerator struct {
	lock sync.Mutex
	id   ClientID
}

func (g *ClientIDGenerator) NextID() ClientID {
	g.lock.Lock()
	out := g.id
	g.id.id++
	g.lock.Unlock()
	return out
}
