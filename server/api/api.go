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

package api

import (
	"github.com/google/uuid"
	"trellis.tech/trellis.v0/service"
)

// TableName default api
var TableName = "api"

// API api struct
type API struct {
	ID             string `xorm:"id"`
	Name           string `xorm:"name"`
	ServiceDomain  string `xorm:"service_domain"`
	ServiceName    string `xorm:"service_name"`
	ServiceVersion string `xorm:"service_version"`
	Topic          string `xorm:"topic"`
	Status         string `xorm:"status"`
}

// TableName database table name
func (*API) TableName() string {
	return TableName
}

func (p *httpServer) syncAPIs(s *service.Service) {
	for {
		syncID := uuid.NewString()
		p.options.Logger.Info("start_sync_apis", "sync", syncID, "service", s)

		params := map[string]interface{}{"`status`": "normal"}

		if s.GetDomain() != "" {
			params["`service_domain`"] = s.GetDomain()
		}

		if s.GetName() != "" {
			params["`service_name`"] = s.GetName()
		}

		if s.GetName() != "" {
			params["`service_version`"] = s.GetVersion()
		}

		var apis []*API
		if err := p.apiEngine.Where(params).Find(&apis); err != nil {
			p.options.Logger.Error("sync_apis_failed", "sync", syncID, "err", err.Error())
			<-p.ticker.C
			continue
		}

		lenAPI := len(apis)
		mapAPIs := make(map[string]*API, lenAPI)

		for i := 0; i < lenAPI; i++ {
			p.options.Logger.Debug("add_new_api", "api", apis[i])
			mapAPIs[apis[i].Name] = apis[i]
		}

		p.syncer.Lock()
		p.apis = mapAPIs
		p.syncer.Unlock()
		p.options.Logger.Info("end_sync_apis", "sync", syncID, "service", s)

		<-p.ticker.C
	}
}
