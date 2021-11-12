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
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"trellis.tech/trellis/common.v0/errcode"
	"trellis.tech/trellis/common.v0/txorm"
	"trellis.tech/trellis/common.v0/types"
	"xorm.io/xorm"

	"trellis.tech/trellis.v0/cmd"
	"trellis.tech/trellis.v0/internal/addr"
	"trellis.tech/trellis.v0/internal/gin_middlewares"
	"trellis.tech/trellis.v0/service"
	"trellis.tech/trellis.v0/service/component"
	"trellis.tech/trellis.v0/service/message"
)

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Name: "trellis-postapi", Version: "v1"},
		NewHTTPServer,
	)
}

var handlers = make(map[string]*gin_middlewares.Handler)

// RegistCustomHandlers register customer's handlers
func RegistCustomHandlers(name, path, method string, fn gin.HandlerFunc) {

	if fn == nil {
		panic("handler function should not be nil")
	}
	_, ok := handlers[name]
	if ok {
		panic(fmt.Errorf("handler isalread exists"))
	}

	handlers[name] = &gin_middlewares.Handler{Name: name, URLPath: path, Method: strings.ToUpper(method), Func: fn}
}

type httpServer struct {
	ginMode string

	mode string // LOCAL, REMOTE

	forwardHeaders []string

	apis map[string]*API

	options component.Options

	srv *http.Server

	ticker    *time.Ticker
	syncer    sync.RWMutex
	apiEngine *xorm.Engine
}

// Response response
type Response struct {
	TraceID   string      `json:"trace_id"`
	TraceIP   string      `json:"trace_ip"`
	Code      uint64      `json:"code"`
	Namespace string      `json:"namespace,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Result    interface{} `json:"result"`
}

// InnerResult result of running component
type InnerResult struct {
	HTTPCode    int
	RedirectURL string
	Body        interface{}
}

// NewHTTPServer new api service
func NewHTTPServer(opts ...component.Option) (component.Component, error) {

	s := &httpServer{
		apis: make(map[string]*API),
	}

	for _, o := range opts {
		o(&s.options)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *httpServer) init() error {

	p.mode = p.options.Config.GetString("mode")
	p.ginMode = p.options.Config.GetString("gin_mode")

	gin.SetMode(p.ginMode)

	apisConf := p.options.Config.GetValuesConfig("apis")

	typ := apisConf.GetString("type", "file")
	switch typ {
	case "file":
		apis := apisConf.GetValuesConfig(typ)

		for _, apiKey := range apis.GetKeys() {
			apiConf := apisConf.GetValuesConfig("file." + apiKey)
			if apiConf == nil {
				return fmt.Errorf("init api failed: %s", apiKey)
			}

			api := &API{
				Name:           apiConf.GetString("api"),
				Topic:          apiConf.GetString("topic"),
				ServiceDomain:  apiConf.GetString("service_domain"),
				ServiceName:    apiConf.GetString("service_name"),
				ServiceVersion: apiConf.GetString("service_version"),
			}

			p.apis[api.Name] = api
		}
	case "mysql":

		databaseConf := apisConf.GetValuesConfig(typ)

		engines, err := txorm.NewEnginesFromConfig(databaseConf)
		if err != nil {
			return err
		}
		p.apiEngine = engines[txorm.DefaultDatabase]

		ticker := types.ParseStringTime(apisConf.GetString("ticker", "30s"))

		p.ticker = time.NewTicker(ticker)

		s := &service.Service{
			Domain:  apisConf.GetString("service_domain"),
			Name:    apisConf.GetString("service_name"),
			Version: apisConf.GetString("service_version"),
		}
		go p.syncAPIs(s)
	default:
		return fmt.Errorf("unknown apis' config type")
	}

	engine := gin.New()

	engine.Use(gin.Recovery(), gin_middlewares.NewRequestID(), gin_middlewares.StatFunc(p.options.Logger))

	httpConf := p.options.Config.GetValuesConfig("http")

	staticPath := httpConf.GetString("static_path")
	staticRedirect := httpConf.GetString("static_redirect")
	staticRoot := httpConf.GetString("static_root", "./root")
	if staticPath != "" {
		if staticRedirect != "" {
			engine.GET(staticPath, func(c *gin.Context) {
				c.Redirect(http.StatusFound, staticRedirect)
			})

			engine.Static(staticRedirect, staticRoot)
		} else {
			engine.Static(staticPath, staticRoot)
		}
	}

	gin_middlewares.LoadPprof(engine, httpConf.GetValuesConfig("pprof"))

	ginHanlders := []gin.HandlerFunc{}

	if gzipH := gin_middlewares.LoadGZip(httpConf.GetValuesConfig("gzip")); gzipH != nil {
		ginHanlders = append(ginHanlders, gzipH)
	}

	if corsH := gin_middlewares.LoadCors(httpConf.GetValuesConfig("cors")); corsH != nil {
		ginHanlders = append(ginHanlders, corsH)
	}

	for _, name := range gin_middlewares.IndexGinFuncs {
		ginHanlders = append(ginHanlders, gin_middlewares.UseFuncs[name])
	}
	engine.Use(ginHanlders...)

	urlPath := httpConf.GetString("postapi")
	if len(urlPath) != 0 {
		engine.POST(urlPath, p.serve)
	}

	for _, v := range handlers {
		p.options.Logger.Info("start_customer_handler", "name", v.Name, "path", v.URLPath, "method", v.Method)
		engine.Handle(v.Method, v.URLPath, v.Func)
	}

	p.forwardHeaders = httpConf.GetStringList("forward.headers")

	p.srv = &http.Server{
		Addr:    httpConf.GetString("address", ":8080"),
		Handler: engine,
	}

	return nil
}

func (p *httpServer) Route(message.Message) (interface{}, error) {
	return nil, nil
}

func (p *httpServer) Start() error {

	go func() {

		var err error

		sslConf := p.options.Config.GetValuesConfig("http.ssl")

		if sslConf != nil && sslConf.GetBoolean("enabled", false) {
			err = p.srv.ListenAndServeTLS(
				sslConf.GetString("cert-file"),
				sslConf.GetString("cert-key"),
			)
		} else {
			err = p.srv.ListenAndServe()
		}

		if err != nil {
			if err != http.ErrServerClosed {
				p.options.Logger.Error("failed_listen_and_serve", "err", err.Error())
				log.Fatalln(err)
			}
		}
	}()
	return nil
}

func (p *httpServer) Stop() error {

	dur := p.options.Config.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errcode.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *httpServer) serve(gCtx *gin.Context) {

	apiName := gCtx.Request.Header.Get(service.HeaderXAPI)
	clientIP := addr.GetClientIP(gCtx.Request)

	reqID := gCtx.GetHeader(service.HeaderXRequestID)

	r := &Response{
		TraceID: reqID,
		TraceIP: addr.ExternalIPs()[0],
	}

	api, ok := p.getAPI(apiName)
	if !ok {
		r.Code = 11
		r.Msg = "api not found"
		r.Namespace = "trellis"
		gCtx.JSON(http.StatusBadRequest, r)
		p.options.Logger.Error("api_not_found", "request_id", reqID, "api_name", apiName, "client_ip", clientIP)
		return
	}

	body, err := gCtx.GetRawData()
	if err != nil {
		r.Code = 10
		r.Msg = fmt.Sprintf("bad request: %s", err.Error())
		r.Namespace = "trellis"
		gCtx.JSON(http.StatusBadRequest, r)
		p.options.Logger.Error("get_raw_data", "request_id", r.TraceID, "api_name", apiName, "client_ip", clientIP, "err", err)
		return
	}

	payload := &message.Payload{
		Header: make(map[string]string),
		Body:   body,
	}

	payload.Set(service.HeaderXClientIP, clientIP)
	payload.Set(service.HeaderXRequestID, reqID)
	for _, h := range p.forwardHeaders {
		payload.Set(h, gCtx.GetHeader(h))
	}

	msg := message.NewMessage(message.Service(
		&service.Service{
			Domain:  api.ServiceDomain,
			Name:    api.ServiceName,
			Version: api.ServiceVersion,
			Topic:   api.Topic}),
		message.MessagePayload(payload),
	)

	resp, err := p.options.Caller.CallComponent(msg)
	if err == nil {
		r.reponse(gCtx, resp)
		return
	}

	// errcode
	switch et := err.(type) {
	case errcode.ErrorCode:
		r.Code = et.Code()
		r.Msg = et.Error()
		r.Namespace = et.Namespace()
	case errcode.SimpleError:
		r.Code = 14
		r.Msg = et.Error()
		r.Namespace = et.Namespace()
	default:
		r.Code = 15
		r.Msg = et.Error()
		r.Namespace = "trellis"
	}

	p.options.Logger.Error("call_server_failed", "request_id", r.TraceID, "api_name", apiName, "client_ip", clientIP, "err", r)
	gCtx.JSON(200, r)
}

func (p *httpServer) getAPI(name string) (*API, bool) {
	p.syncer.RLock()
	api, ok := p.apis[name]
	p.syncer.RUnlock()
	return api, ok
}

func (p *Response) reponse(ctx *gin.Context, resp interface{}) {
	switch t := resp.(type) {
	case InnerResult:
		p.genResult(ctx, &t)
		return
	case *InnerResult:
		p.genResult(ctx, t)
		return
	default:
		p.Result = resp
		ctx.JSON(200, p)
		return
	}
}

func (p *Response) genResult(ctx *gin.Context, t *InnerResult) {
	switch t.HTTPCode {
	case http.StatusMovedPermanently, http.StatusFound:
		ctx.Redirect(t.HTTPCode, t.RedirectURL)
	default:
		p.Result = t.Body
		ctx.JSON(t.HTTPCode, p)
	}
}
