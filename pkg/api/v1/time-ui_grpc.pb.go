// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// TimeUIClient is the client API for TimeUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimeUIClient interface {
	// ListEngineSpecs returns a list of Time Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (TimeUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type timeUIClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeUIClient(cc grpc.ClientConnInterface) TimeUIClient {
	return &timeUIClient{cc}
}

func (c *timeUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (TimeUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimeUI_ServiceDesc.Streams[0], "/v1.TimeUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &timeUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TimeUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type timeUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *timeUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *timeUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.TimeUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimeUIServer is the server API for TimeUI service.
// All implementations must embed UnimplementedTimeUIServer
// for forward compatibility
type TimeUIServer interface {
	// ListEngineSpecs returns a list of Time Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, TimeUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedTimeUIServer()
}

// UnimplementedTimeUIServer must be embedded to have forward compatible implementations.
type UnimplementedTimeUIServer struct {
}

func (UnimplementedTimeUIServer) ListEngineSpecs(*ListEngineSpecsRequest, TimeUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedTimeUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedTimeUIServer) mustEmbedUnimplementedTimeUIServer() {}

// UnsafeTimeUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimeUIServer will
// result in compilation errors.
type UnsafeTimeUIServer interface {
	mustEmbedUnimplementedTimeUIServer()
}

func RegisterTimeUIServer(s grpc.ServiceRegistrar, srv TimeUIServer) {
	s.RegisterService(&TimeUI_ServiceDesc, srv)
}

func _TimeUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TimeUIServer).ListEngineSpecs(m, &timeUIListEngineSpecsServer{stream})
}

type TimeUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type timeUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *timeUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _TimeUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TimeUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TimeUI_ServiceDesc is the grpc.ServiceDesc for TimeUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimeUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.TimeUI",
	HandlerType: (*TimeUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _TimeUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _TimeUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "time-ui.proto",
}