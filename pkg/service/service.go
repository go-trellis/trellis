package service

import (
	"strings"
)

func (m *Service) FullPath() string {
	ss := []string{ReplaceURL(m.GetDomain()), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion())}
	return strings.Join(ss, "/")
}

func (m *Service) GetPath(registry string) string {
	ss := []string{registry, ReplaceURL(m.GetDomain()), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion())}
	return strings.Join(ss, "/")
}

func (m *ServiceNode) RegisteredServiceNode(registry string) string {
	ss := []string{registry,
		ReplaceURL(m.GetService().GetDomain()), ReplaceURL(m.GetService().GetName()), ReplaceURL(m.GetService().GetVersion()),
		ReplaceURL(m.GetNode().GetValue())}
	return strings.Join(ss, "/")
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
