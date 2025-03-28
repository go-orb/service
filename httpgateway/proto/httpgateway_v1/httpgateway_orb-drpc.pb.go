// Code generated by protoc-gen-go-orb. DO NOT EDIT.
//
// version:
// - protoc-gen-go-orb        v0.0.1
// - protoc                   v5.29.2
//
// source: httpgateway_v1/httpgateway.proto

package httpgateway_v1

import (
	context "context"
	errors "errors"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_httpgateway_v1_httpgateway_proto struct{}

func (drpcEncoding_File_httpgateway_v1_httpgateway_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_httpgateway_v1_httpgateway_proto) MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(buf, msg.(proto.Message))
}

func (drpcEncoding_File_httpgateway_v1_httpgateway_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_httpgateway_v1_httpgateway_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_httpgateway_v1_httpgateway_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}

type DRPCHttpGatewayServer interface {
	AddRoutes(context.Context, *Routes) (*emptypb.Empty, error)
	SetRoutes(context.Context, *Routes) (*emptypb.Empty, error)
	RemoveRoutes(context.Context, *Paths) (*emptypb.Empty, error)
}

type DRPCHttpGatewayUnimplementedServer struct{}

func (s *DRPCHttpGatewayUnimplementedServer) AddRoutes(context.Context, *Routes) (*emptypb.Empty, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCHttpGatewayUnimplementedServer) SetRoutes(context.Context, *Routes) (*emptypb.Empty, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCHttpGatewayUnimplementedServer) RemoveRoutes(context.Context, *Paths) (*emptypb.Empty, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCHttpGatewayDescription struct{}

func (DRPCHttpGatewayDescription) NumMethods() int { return 3 }

func (DRPCHttpGatewayDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/httpgateway.v1.HttpGateway/AddRoutes", drpcEncoding_File_httpgateway_v1_httpgateway_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCHttpGatewayServer).
					AddRoutes(
						ctx,
						in1.(*Routes),
					)
			}, DRPCHttpGatewayServer.AddRoutes, true
	case 1:
		return "/httpgateway.v1.HttpGateway/SetRoutes", drpcEncoding_File_httpgateway_v1_httpgateway_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCHttpGatewayServer).
					SetRoutes(
						ctx,
						in1.(*Routes),
					)
			}, DRPCHttpGatewayServer.SetRoutes, true
	case 2:
		return "/httpgateway.v1.HttpGateway/RemoveRoutes", drpcEncoding_File_httpgateway_v1_httpgateway_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCHttpGatewayServer).
					RemoveRoutes(
						ctx,
						in1.(*Paths),
					)
			}, DRPCHttpGatewayServer.RemoveRoutes, true
	default:
		return "", nil, nil, nil, false
	}
}

type DRPCHttpGateway_AddRoutesStream interface {
	drpc.Stream
	SendAndClose(*emptypb.Empty) error
}

type drpcHttpGateway_AddRoutesStream struct {
	drpc.Stream
}

func (x *drpcHttpGateway_AddRoutesStream) GetStream() drpc.Stream {
	return x.Stream
}

func (x *drpcHttpGateway_AddRoutesStream) SendAndClose(m *emptypb.Empty) error {
	if err := x.MsgSend(m, drpcEncoding_File_httpgateway_v1_httpgateway_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCHttpGateway_SetRoutesStream interface {
	drpc.Stream
	SendAndClose(*emptypb.Empty) error
}

type drpcHttpGateway_SetRoutesStream struct {
	drpc.Stream
}

func (x *drpcHttpGateway_SetRoutesStream) GetStream() drpc.Stream {
	return x.Stream
}

func (x *drpcHttpGateway_SetRoutesStream) SendAndClose(m *emptypb.Empty) error {
	if err := x.MsgSend(m, drpcEncoding_File_httpgateway_v1_httpgateway_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCHttpGateway_RemoveRoutesStream interface {
	drpc.Stream
	SendAndClose(*emptypb.Empty) error
}

type drpcHttpGateway_RemoveRoutesStream struct {
	drpc.Stream
}

func (x *drpcHttpGateway_RemoveRoutesStream) GetStream() drpc.Stream {
	return x.Stream
}

func (x *drpcHttpGateway_RemoveRoutesStream) SendAndClose(m *emptypb.Empty) error {
	if err := x.MsgSend(m, drpcEncoding_File_httpgateway_v1_httpgateway_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
