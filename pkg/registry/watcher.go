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

package registry

import (
	"encoding/json"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/service"
)

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

type WatchService struct {
	Service  *service.Service `yaml:"service" json:"service"`
	NodeType node.NodeType    `yaml:"node_type" json:"node_type"`

	Metadata *WatchServiceMetadata `yaml:"metadata" json:"metadata"`
}

type WatchServiceMetadata struct {
	ClientConfig *clients.Config `yaml:"client_config" json:"client_config"`
}

func ToWatchServiceMetadata(data string) (*WatchServiceMetadata, error) {
	m := &WatchServiceMetadata{}
	err := json.Unmarshal([]byte(data), m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
