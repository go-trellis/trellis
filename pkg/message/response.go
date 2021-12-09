package message

import (
	"reflect"

	"trellis.tech/trellis.v1/pkg/codec"

	"trellis.tech/trellis/common.v0/json"
)

type Option func(*Options)

type Options struct {
	Code int64
	Err  error
}

func Code(code int64) Option {
	return func(o *Options) {
		o.Code = code
	}
}

func Error(err error) Option {
	return func(o *Options) {
		o.Err = err
	}
}

func NewResponse(data interface{}, opts ...Option) *Response {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	resp := &Response{
		Code: options.Code,
		Payload: &Payload{
			Header: make(map[string]string),
		},
	}

	if options.Err != nil {
		resp.ErrMsg = options.Err.Error()
	}

	if data == nil {
		return resp
	}

	var respData interface{}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr, reflect.Struct:
		respData = data
	default:
		respData = map[string]interface{}{"data": data}
	}
	bs, err := json.Marshal(respData)
	if err != nil {
		resp.Code = 500
		resp.ErrMsg = err.Error()
		return resp
	}
	resp.Payload.Header["Content-Type"] = codec.ContentTypeJson
	resp.Payload.Body = bs
	return resp
}
