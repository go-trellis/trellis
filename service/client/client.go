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

package client

import (
	"context"

	"trellis.tech/trellis.v0/service"
	"trellis.tech/trellis.v0/service/codec"
)

// DefaultClient implementation
var DefaultClient Client

// Client client
type Client interface {
	NewMessage(msg interface{}, opts ...MessageOption) Message
	NewRequest(service *service.Service, endpoint string, req interface{}, reqOpts ...RequestOption) Request

	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error

	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)

	String() string
}

// Message is the interface for publishing asynchronously
type Message interface {
	Topic() string
	Payload() interface{}
	ContentType() string
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	// Service The service to call
	Service() *service.Service
	// Method The action to take
	Method() string
	// Endpoint The endpoint to invoke
	Endpoint() string
	// ContentType The content type
	ContentType() string
	// Body The unencoded request body
	Body() interface{}
	// Codec Write to the encoded request writer. This is nil before a call is made
	Codec() codec.Codec
	// Stream indicates whether the request will be a streaming one rather than unary
	Stream() bool
}

// Stream is the interface for a bidirectional synchronous stream
type Stream interface {
	// Context for the stream
	Context() context.Context
	// Request The request made
	Request() Request
	// Response The response read
	Response() Response
	// Send will encode and send a request
	Send(interface{}) error
	// Recv will decode and read a response
	Recv(interface{}) error
	// Error returns the stream error
	Error() error
	// Close closes the stream
	Close() error
}

// Response is the response received from a service
type Response interface {
	// Codec Reader the response
	Codec() codec.Codec
	// Header read the header
	Header() map[string]string
	// Read the undecoded response
	Read() ([]byte, error)
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

// RequestOption used by NewRequest
type RequestOption func(*RequestOptions)

// MessageOption used by NewMessage
type MessageOption func(*MessageOptions)
