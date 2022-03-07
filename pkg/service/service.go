package service

import (
	"path/filepath"
	"strings"

	"trellis.tech/trellis.v1/pkg/node"
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
	if m == nil {
		return ""
	}
	m.Domain = checkDomain(m.GetDomain())
	return filepath.Join(m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()))
}

func (m *Service) TopicPath() string {
	if m == nil {
		return ""
	}
	m.Domain = checkDomain(m.GetDomain())
	return filepath.Join(m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()), ReplaceURL(m.GetTopic()))
}

func (m *Service) GetPath(registry string) string {
	m.Domain = checkDomain(m.GetDomain())
	return filepath.Join(registry, m.GetDomain(), ReplaceURL(m.GetName()), ReplaceURL(m.GetVersion()))
}

type Node struct {
	Service *Service
	Node    *node.Node
}

func (m *Node) RegisteredServiceNode(registry string) string {
	m.Service.Domain = checkDomain(m.Service.GetDomain())
	return filepath.Join(registry, m.Service.GetDomain(),
		ReplaceURL(m.Service.GetName()),
		ReplaceURL(m.Service.GetVersion()),
		ReplaceURL(m.Node.GetValue()))
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
