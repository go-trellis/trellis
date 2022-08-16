/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package memory

import (
	"testing"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/logger"
)

func TestWatcher(t *testing.T) {
	w := &Watcher{
		id:     "test",
		res:    make(chan *registry.Result),
		exit:   make(chan bool),
		serv:   &service.Service{},
		logger: logger.Noop(),
	}

	go func() {
		w.res <- &registry.Result{}
	}()

	_, err := w.Next()
	if err != nil {
		t.Fatal("unexpected err", err)
	}

	w.Stop()

	if _, err := w.Next(); err != nil {
		t.Fatal("expected error on Next()")
	}
}
