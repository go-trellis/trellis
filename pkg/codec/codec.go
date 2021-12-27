package codec

import (
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
)

type NewCodecFunc func() (Codec, error)

var (
	codecs        map[string]NewCodecFunc = make(map[string]NewCodecFunc)
	defaultCodecs map[string]Codec        = make(map[string]Codec)
)

func init() {
	RegisterCodec(mime.ContentTypeJson, NewJsonCodec)

	jsonCodec, _ := NewCodec(mime.ContentTypeJson)
	defaultCodecs[mime.ContentTypeJson] = jsonCodec
}

type Codec interface {
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}

func NewCodec(name string) (Codec, error) {
	fn, exist := codecs[name]
	if !exist {
		return nil, errcode.Newf("codec not exist %q", name)
	}
	return fn()
}

func RegisterCodec(name string, fn NewCodecFunc) {
	if len(name) == 0 {
		panic("codec name is empty")
	}

	if fn == nil {
		panic("codec fn is nil")
	}

	c, err := fn()
	if err != nil {
		panic(err)
	}

	defaultCodecs[name] = c
	codecs[name] = fn
}

func Select(contentType string) Codec {
	return defaultCodecs[contentType]
}
