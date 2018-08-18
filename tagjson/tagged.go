// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tagjson

import (
	"encoding/json"
	"errors"
	"fmt"
)

var errNotTagged = errors.New("not proper tagged message")

func TagWithType(v interface{}) ([]byte, error) {
	tag := fmt.Sprintf("%T", v)
	tagged := map[string]interface{}{tag: v}
	return json.Marshal(tagged)
}

func Untag(rawMsg json.RawMessage) (tag string, msg json.RawMessage, err error) {
	tagged := make(map[string]json.RawMessage, 1)
	err = json.Unmarshal(rawMsg, &tagged)
	if err != nil {
		return "", []byte{}, err
	}
	if len(tagged) != 1 {
		return "", []byte{}, errNotTagged
	}
	for tag, msg = range tagged {
		return tag, msg, nil
	}
	panic("unreachable")
}
