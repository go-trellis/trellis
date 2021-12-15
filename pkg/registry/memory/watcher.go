package memory

import (
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v0/logger"
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
