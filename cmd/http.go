package cmd

import (
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/server/http"
	"trellis.tech/trellis.v1/pkg/trellis"
)

func NewHTTPServer(cfg *trellis.HTTPServerConfig, r router.Router) (server.Server, error) {

	s, err := http.NewServer(http.Config(cfg), http.Router(r))
	if err != nil {
		return nil, err
	}

	var hs []*http.Handler
	for _, hCfg := range cfg.Handlers {
		h, err := getHTTPHandler(hCfg)
		if err != nil {
			return nil, err
		}
		hs = append(hs, h)
	}
	s.RegisterHandler(hs...)

	for _, group := range cfg.Groups {
		hgs, err := getHTTPGroupHandlers(group)
		if err != nil {
			return nil, err
		}
		s.RegisterGroup(group.Path, hgs...)
	}

	return s, err
}

func getHTTPGroupHandlers(group *trellis.HTTPGroup) ([]*http.Handler, error) {

	var hs []*http.Handler
	for _, hCfg := range group.Handlers {
		h, err := getHTTPHandler(hCfg)
		if err != nil {
			return nil, err
		}
		hs = append(hs, h)
	}

	return hs, nil
}

func getHTTPHandler(handler *trellis.HTTPHandler) (*http.Handler, error) {
	h := &http.Handler{
		Method: handler.Method,
		Path:   handler.Path,
	}

	for _, use := range handler.Uses {
		uFunc, err := server.GetUseFunc(use)
		if err != nil {
			return nil, err
		}
		h.Uses = append(h.Uses, uFunc)
	}

	if handler.Handler != "" {
		uFunc, err := server.GetUseFunc(handler.Handler)
		if err != nil {
			return nil, err
		}
		h.Uses = append(h.Uses, uFunc)
	}
	return h, nil
}
