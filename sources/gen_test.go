// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sources_test

//go:generate mockgen -destination ./mock_io_test.go -package sources_test io ReadCloser,WriteCloser
//go:generate mockgen -destination ./mock_factory_test.go -package sources_test github.com/szabba/munch/sources InputFactory,ParserFactory
