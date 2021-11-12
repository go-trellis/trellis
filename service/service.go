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

package service

import (
	"errors"
	"strings"

	"trellis.tech/trellis/common.v0/crypto/hash"
)

const (
	defDomain = "trellis"
	registry  = "/trellis/registry"
)

func (m *Service) init() {
	if m.Domain == "" {
		m.Domain = defDomain
	}
}

// ID gen service id
func (m *Service) ID(ps ...string) string {
	if m == nil {
		return ""
	}
	m.init()
	return hash.NewCRCIEEE().Sum(m.TrellisPath(ps...))
}

// TrellisName service name
func (m *Service) TrellisName() string {
	if m == nil {
		return ""
	}
	m.init()
	return joinPath([]string{ReplaceURL(m.Domain), ReplaceURL(m.Name)})
}

// TrellisPath Service full path
func (m *Service) TrellisPath(ps ...string) string {
	if m == nil {
		return ""
	}
	m.init()

	ss := []string{ReplaceURL(m.Domain), ReplaceURL(m.Name), ReplaceURL(m.Version)}

	for _, s := range ps {
		ss = append(ss, ReplaceURL(s))
	}

	return joinPath(ss)
}

// FullRegistryPath Service full registry path
func (m *Service) FullRegistryPath(ps ...string) string {
	if m == nil {
		return ""
	}

	m.init()

	ss := []string{registry, ReplaceURL(m.Domain), ReplaceURL(m.Name), ReplaceURL(m.Version)}

	for _, s := range ps {
		ss = append(ss, ReplaceURL(s))
	}

	return joinPath(ss)
}

func joinPath(ss []string) string {
	return strings.Join(ss, "/")
}

// ParseService parse a string to base service
func ParseService(s string) (*Service, error) {
	ss := strings.Split(s, "/")

	lenSS := len(ss)

	var bs *Service
	if lenSS == 3 {
		bs = &Service{
			Name:    ss[1],
			Version: ss[2],
		}
	} else if lenSS > 3 {
		bs = &Service{
			Domain:  ss[1],
			Name:    ss[2],
			Version: ss[3],
		}
	} else {
		return nil, errors.New("failed parse base service")
	}

	bs.init()

	return bs, nil
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
