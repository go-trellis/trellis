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

package service

import (
	"path/filepath"
	"strings"
)

const DefaultDomain = "/trellis"

func NewService(domain, name, version string) *Service {
	domain = CheckDomain(domain)
	return &Service{Domain: domain, Name: name, Version: version}
}

func NewServiceWithTopic(domain, name, version, topic string) *Service {
	s := NewService(domain, name, version)
	s.Topic = topic
	return s
}

func CheckDomain(domain string) string {
	domain = ReplaceURL(strings.TrimLeft(domain, "/"))
	if domain == "" {
		domain = DefaultDomain
	} else if !strings.HasPrefix(domain, "/") {
		domain = "/" + domain
	}
	return domain
}

func (m *Service) FullPath() string {
	if m == nil {
		return ""
	}
	m.Domain = CheckDomain(m.GetDomain())
	return filepath.Join(m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()))
}

func (m *Service) TopicPath() string {
	if m == nil {
		return ""
	}
	m.Domain = CheckDomain(m.GetDomain())
	return filepath.Join(m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()), ReplaceURL(m.GetTopic()))
}

func (m *Service) GetPath(registry string) string {
	m.Domain = CheckDomain(m.GetDomain())
	return filepath.Join(registry, m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()))
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
