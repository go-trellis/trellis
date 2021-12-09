/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

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
	"trellis.tech/trellis/common.v0/logger"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"
)

// Watcher watcher
type Watcher struct {
	id   string
	exit chan bool
	res  chan *registry.Result

	logger logger.Logger

	prefix string
	serv   *service.Service
}

// Next watch the register result
func (p *Watcher) Next() (*registry.Result, error) {
	for {
		select {
		case r := <-p.res:
			if p.serv != nil && p.serv.Name != "" &&
				p.serv.GetPath(p.prefix) != r.ServiceNode.GetService().GetPath(p.prefix) {
				continue
			}

			p.logger.Infof("watcher next", "id", p.id, "prefix", p.prefix, "service_node", r.ServiceNode)
			return r, nil
		case <-p.exit:
			p.logger.Infof("watcher stop", "id", p.id, "prefix", p.prefix, "service", p.serv)
			return nil, nil
		}
	}
}

// Stop watcher
func (p *Watcher) Stop() {
	select {
	case <-p.exit:
		return
	default:
		close(p.exit)
	}
}
