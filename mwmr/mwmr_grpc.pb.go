// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: mwmr/mwmr.proto

package mwmr

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MWMRClient is the client API for MWMR service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MWMRClient interface {
	GetPhase(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	SetPhase(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetACK, error)
}

type mWMRClient struct {
	cc grpc.ClientConnInterface
}

func NewMWMRClient(cc grpc.ClientConnInterface) MWMRClient {
	return &mWMRClient{cc}
}

func (c *mWMRClient) GetPhase(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/mwmr.MWMR/GetPhase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mWMRClient) SetPhase(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetACK, error) {
	out := new(SetACK)
	err := c.cc.Invoke(ctx, "/mwmr.MWMR/SetPhase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MWMRServer is the server API for MWMR service.
// All implementations must embed UnimplementedMWMRServer
// for forward compatibility
type MWMRServer interface {
	GetPhase(context.Context, *GetRequest) (*GetReply, error)
	SetPhase(context.Context, *SetRequest) (*SetACK, error)
	mustEmbedUnimplementedMWMRServer()
}

// UnimplementedMWMRServer must be embedded to have forward compatible implementations.
type UnimplementedMWMRServer struct {
}

func (UnimplementedMWMRServer) GetPhase(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhase not implemented")
}
func (UnimplementedMWMRServer) SetPhase(context.Context, *SetRequest) (*SetACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPhase not implemented")
}
func (UnimplementedMWMRServer) mustEmbedUnimplementedMWMRServer() {}

// UnsafeMWMRServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MWMRServer will
// result in compilation errors.
type UnsafeMWMRServer interface {
	mustEmbedUnimplementedMWMRServer()
}

func RegisterMWMRServer(s grpc.ServiceRegistrar, srv MWMRServer) {
	s.RegisterService(&MWMR_ServiceDesc, srv)
}

func _MWMR_GetPhase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MWMRServer).GetPhase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mwmr.MWMR/GetPhase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MWMRServer).GetPhase(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MWMR_SetPhase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MWMRServer).SetPhase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mwmr.MWMR/SetPhase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MWMRServer).SetPhase(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MWMR_ServiceDesc is the grpc.ServiceDesc for MWMR service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MWMR_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mwmr.MWMR",
	HandlerType: (*MWMRServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPhase",
			Handler:    _MWMR_GetPhase_Handler,
		},
		{
			MethodName: "SetPhase",
			Handler:    _MWMR_SetPhase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mwmr/mwmr.proto",
}
