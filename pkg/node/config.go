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

package node

import (
	"fmt"

	"trellis.tech/trellis/common.v1/config"
)

// NewNodesFromConfig 同步配置文件
func NewNodesFromConfig(filepath string) (map[string]Manager, error) {
	cfg, err := config.NewConfigOptions(config.OptionFile(filepath))
	if err != nil {
		return nil, err
	}
	return NewNodes(cfg)
}

// NewNodes 增加Nodes节点
func NewNodes(cfg config.Config) (ms map[string]Manager, err error) {
	mapManager := make(map[string]Manager)

	valConfigs := cfg.GetValuesConfig("node")
	for _, key := range valConfigs.GetKeys() {
		m, err := New(NodeType(valConfigs.GetInt(key+".type")), key)
		if err != nil {
			return nil, err
		}
		nodesCfg := valConfigs.GetValuesConfig(key + ".nodes")

		for _, nKey := range nodesCfg.GetKeys() {
			item := &Node{}
			item.Value = nodesCfg.GetString(nKey + ".value")
			item.Weight = uint32(nodesCfg.GetInt(nKey + ".weight"))
			for k, v := range nodesCfg.GetMap(nKey + ".metadata") {
				if v == nil {
					continue
				}
				switch t := v.(type) {
				case int:
					item.Metadata[k] = fmt.Sprintf("%d", t)
				case float64:
					item.Metadata[k] = fmt.Sprintf("%f", t)
				case string:
					item.Metadata[k] = t
				default:
					continue
				}
			}
			m.Add(item)
		}
		mapManager[key] = m
	}
	return mapManager, nil
}
