package local

import (
	"context"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis/common.v0/errcode"
	"trellis.tech/trellis/common.v0/json"
)

var (
	_ clients.Client = (*Client)(nil)

	c = &Client{}
)

type Client struct{}

func NewClient() (clients.Client, error) {
	return c, nil
}

func (*Client) Call(_ context.Context, in *message.Request, _ ...clients.CallOption) (*message.Response, error) {
	comp := router.GetComponent(in.GetService())
	if comp == nil {
		return nil, errcode.Newf("not found component: %s", in.GetService().FullPath())
	}
	hResp, err := comp.Route(in.GetService().GetTopic(), in.GetPayload())
	if err != nil {
		// TODO log err
		return nil, err
	}

	if hResp == nil {
		return &message.Response{
			Code: 0,
		}, nil
	}

	switch t := hResp.(type) {
	case message.Response:
		return &t, nil
	case *message.Response:
		return t, nil
	case *message.Payload:
		return &message.Response{
			Code:    0,
			Payload: t,
		}, nil
	case message.Payload:
		return &message.Response{
			Code:    0,
			Payload: &t,
		}, nil
	default:
		bs, err := json.Marshal(hResp)
		if err != nil {
			return nil, err
		}
		return &message.Response{
			Code: 0,
			Payload: &message.Payload{
				Header: map[string]string{"Content-Type": codec.ContentTypeJson},
				Body:   bs,
			},
		}, nil
	}
}