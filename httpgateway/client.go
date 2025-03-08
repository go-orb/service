package service

import (
	"context"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/service/httpgateway/proto/httpgateway_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	serviceName string
	clientWire  client.Type
}

func New(serviceName string, clientWire client.Type) *Client {
	return &Client{
		serviceName: serviceName,
		clientWire:  clientWire,
	}
}

func (c *Client) ServiceName() string {
	return c.serviceName
}

func (c *Client) AddRoutes(ctx context.Context, routes *httpgateway_v1.Routes) (*emptypb.Empty, error) {
	cli := httpgateway_v1.NewHttpGatewayClient(c.clientWire)

	return cli.AddRoutes(ctx, c.serviceName, routes)
}

func (c *Client) SetRoutes(ctx context.Context, routes *httpgateway_v1.Routes) (*emptypb.Empty, error) {
	cli := httpgateway_v1.NewHttpGatewayClient(c.clientWire)

	return cli.SetRoutes(ctx, c.serviceName, routes)
}

func (c *Client) RemoveRoutes(ctx context.Context, paths *httpgateway_v1.Paths) (*emptypb.Empty, error) {
	cli := httpgateway_v1.NewHttpGatewayClient(c.clientWire)

	return cli.RemoveRoutes(ctx, c.serviceName, paths)
}
