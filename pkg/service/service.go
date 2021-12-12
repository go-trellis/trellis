package service

import (
	"strings"
)

const DefaultDomain = "/trellis"

func NewService(domain, name, version string) *Service {
	domain = checkDomain(domain)
	return &Service{Domain: domain, Name: name, Version: version}
}

func checkDomain(domain string) string {
	domain = ReplaceURL(strings.TrimLeft(domain, "/"))
	if domain == "" {
		domain = DefaultDomain
	} else if !strings.HasPrefix(domain, "/") {
		domain = "/" + domain
	}
	return domain
}

func (m *Service) FullPath() string {
	m.Domain = checkDomain(m.GetDomain())
	ss := []string{m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion())}
	return strings.Join(ss, "/")
}

func (m *Service) GetPath(registry string) string {
	m.Domain = checkDomain(m.GetDomain())
	ss := []string{registry, m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion())}
	return strings.Join(ss, "/")
}

func (m *ServiceNode) RegisteredServiceNode(registry string) string {
	m.Service.Domain = checkDomain(m.GetService().GetDomain())
	ss := []string{registry, m.GetService().GetDomain(),
		ReplaceURL(m.GetService().GetName()), ReplaceURL(m.GetService().GetVersion()), ReplaceURL(m.GetNode().GetValue())}
	return strings.Join(ss, "/")
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
