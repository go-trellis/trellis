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

package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"trellis.tech/trellis.v0/cmd"
	"trellis.tech/trellis.v0/internal/addr"
	"trellis.tech/trellis.v0/internal/gin_middlewares"
	"trellis.tech/trellis.v0/server"
	"trellis.tech/trellis.v0/service"
	"trellis.tech/trellis.v0/service/component"
	"trellis.tech/trellis.v0/service/message"
	"trellis.tech/trellis/common.v0/errcode"
)

var s = &service.Service{Name: "trellis-server-http", Version: "v1"}

func init() {
	_ = cmd.DefaultCompManager.RegisterComponentFunc(s, NewHTTPServer)
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

	serverIP string

	forwardHeaders []string

	options component.Options

	srv *http.Server
}

// NewHTTPServer new api service
func NewHTTPServer(opts ...component.Option) (component.Component, error) {
	s := &httpServer{}
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

	ips := addr.ExternalIPs()
	if len(ips) > 0 {
		p.serverIP = ips[0]
	} else {
		p.serverIP = "unknown server ip"
	}

	p.ginMode = p.options.Config.GetString("gin_mode")

	gin.SetMode(p.ginMode)

	engine := gin.New()

	engine.Use(gin.Recovery(), gin_middlewares.NewRequestID(), gin_middlewares.StatFunc(p.options.Logger))

	httpConf := p.options.Config.GetValuesConfig("http")

	gin_middlewares.LoadPprof(engine, httpConf.GetValuesConfig("pprof"))

	var ginHandlers []gin.HandlerFunc

	if gzipH := gin_middlewares.LoadGZip(httpConf.GetValuesConfig("gzip")); gzipH != nil {
		ginHandlers = append(ginHandlers, gzipH)
	}

	ginHandlers = append(ginHandlers, gin_middlewares.LoadCors(httpConf.GetValuesConfig("cors")))

	for _, name := range gin_middlewares.IndexGinFuncs {
		ginHandlers = append(ginHandlers, gin_middlewares.UseFuncs[name])
	}
	engine.Use(ginHandlers...)

	urlPath := httpConf.GetString("postapi")
	if len(urlPath) != 0 {
		engine.POST(urlPath, p.serve)
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

func (p *httpServer) serve(ctx *gin.Context) {

	clientIP := addr.GetClientIP(ctx.Request)

	reqID := ctx.GetHeader(service.HeaderXRequestID)

	r := &server.Response{
		RequestID: reqID,
		ClientIP:  clientIP,
		ServerIP:  p.serverIP,
	}

	remoteMsg := &message.RemoteMessage{}

	if err := ctx.BindJSON(remoteMsg); err != nil {
		r.Code = 10
		r.Msg = fmt.Sprintf("bad request: %s", err.Error())
		r.Namespace = s.TrellisPath()
		ctx.JSON(http.StatusBadRequest, r)
		p.options.Logger.Error("get_raw_data", "err", err)
		return
	}

	msg := remoteMsg.ToMessage()

	resp, err := p.options.Caller.CallComponent(msg)
	if err == nil {
		switch t := resp.(type) {
		case server.Response:
			r = &t
		case *server.Response:
			r = t
		default:
			r.Result = resp
		}

		ctx.JSON(200, r)
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
		r.Namespace = s.TrellisPath()
	}

	ctx.JSON(200, r)
}
