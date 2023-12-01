// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

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

// FollowerClient is the client API for Follower service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FollowerClient interface {
	SendExecuteCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Nothing, error)
	SendHeartBeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	SendElectionRequest(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionResponse, error)
}

type followerClient struct {
	cc grpc.ClientConnInterface
}

func NewFollowerClient(cc grpc.ClientConnInterface) FollowerClient {
	return &followerClient{cc}
}

func (c *followerClient) SendExecuteCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/protocol.Follower/SendExecuteCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerClient) SendHeartBeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/protocol.Follower/SendHeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerClient) SendElectionRequest(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionResponse, error) {
	out := new(ElectionResponse)
	err := c.cc.Invoke(ctx, "/protocol.Follower/SendElectionRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FollowerServer is the server API for Follower service.
// All implementations must embed UnimplementedFollowerServer
// for forward compatibility
type FollowerServer interface {
	SendExecuteCommand(context.Context, *Command) (*Nothing, error)
	SendHeartBeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	SendElectionRequest(context.Context, *ElectionRequest) (*ElectionResponse, error)
	mustEmbedUnimplementedFollowerServer()
}

// UnimplementedFollowerServer must be embedded to have forward compatible implementations.
type UnimplementedFollowerServer struct {
}

func (UnimplementedFollowerServer) SendExecuteCommand(context.Context, *Command) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendExecuteCommand not implemented")
}
func (UnimplementedFollowerServer) SendHeartBeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartBeat not implemented")
}
func (UnimplementedFollowerServer) SendElectionRequest(context.Context, *ElectionRequest) (*ElectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendElectionRequest not implemented")
}
func (UnimplementedFollowerServer) mustEmbedUnimplementedFollowerServer() {}

// UnsafeFollowerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FollowerServer will
// result in compilation errors.
type UnsafeFollowerServer interface {
	mustEmbedUnimplementedFollowerServer()
}

func RegisterFollowerServer(s grpc.ServiceRegistrar, srv FollowerServer) {
	s.RegisterService(&Follower_ServiceDesc, srv)
}

func _Follower_SendExecuteCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).SendExecuteCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Follower/SendExecuteCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).SendExecuteCommand(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _Follower_SendHeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).SendHeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Follower/SendHeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).SendHeartBeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Follower_SendElectionRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ElectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).SendElectionRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Follower/SendElectionRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).SendElectionRequest(ctx, req.(*ElectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Follower_ServiceDesc is the grpc.ServiceDesc for Follower service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Follower_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Follower",
	HandlerType: (*FollowerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendExecuteCommand",
			Handler:    _Follower_SendExecuteCommand_Handler,
		},
		{
			MethodName: "SendHeartBeat",
			Handler:    _Follower_SendHeartBeat_Handler,
		},
		{
			MethodName: "SendElectionRequest",
			Handler:    _Follower_SendElectionRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}
