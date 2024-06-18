// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: return.proto

package alpha_pos_system_sales_service

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

// PosReturnServiceClient is the client API for PosReturnService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PosReturnServiceClient interface {
	CreatePosReturn(ctx context.Context, in *CreatePosReturnRequest, opts ...grpc.CallOption) (*CreatePosReturnResponse, error)
	ReadPosReturn(ctx context.Context, in *ReadPosReturnRequest, opts ...grpc.CallOption) (*ReadPosReturnResponse, error)
	UpdatePosReturn(ctx context.Context, in *UpdatePosReturnRequest, opts ...grpc.CallOption) (*UpdatePosReturnResponse, error)
	DeletePosReturn(ctx context.Context, in *DeletePosReturnRequest, opts ...grpc.CallOption) (*DeletePosReturnResponse, error)
	ReadAllPosReturns(ctx context.Context, in *ReadAllPosReturnsRequest, opts ...grpc.CallOption) (*ReadAllPosReturnsResponse, error)
}

type posReturnServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPosReturnServiceClient(cc grpc.ClientConnInterface) PosReturnServiceClient {
	return &posReturnServiceClient{cc}
}

func (c *posReturnServiceClient) CreatePosReturn(ctx context.Context, in *CreatePosReturnRequest, opts ...grpc.CallOption) (*CreatePosReturnResponse, error) {
	out := new(CreatePosReturnResponse)
	err := c.cc.Invoke(ctx, "/pos.PosReturnService/CreatePosReturn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posReturnServiceClient) ReadPosReturn(ctx context.Context, in *ReadPosReturnRequest, opts ...grpc.CallOption) (*ReadPosReturnResponse, error) {
	out := new(ReadPosReturnResponse)
	err := c.cc.Invoke(ctx, "/pos.PosReturnService/ReadPosReturn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posReturnServiceClient) UpdatePosReturn(ctx context.Context, in *UpdatePosReturnRequest, opts ...grpc.CallOption) (*UpdatePosReturnResponse, error) {
	out := new(UpdatePosReturnResponse)
	err := c.cc.Invoke(ctx, "/pos.PosReturnService/UpdatePosReturn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posReturnServiceClient) DeletePosReturn(ctx context.Context, in *DeletePosReturnRequest, opts ...grpc.CallOption) (*DeletePosReturnResponse, error) {
	out := new(DeletePosReturnResponse)
	err := c.cc.Invoke(ctx, "/pos.PosReturnService/DeletePosReturn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posReturnServiceClient) ReadAllPosReturns(ctx context.Context, in *ReadAllPosReturnsRequest, opts ...grpc.CallOption) (*ReadAllPosReturnsResponse, error) {
	out := new(ReadAllPosReturnsResponse)
	err := c.cc.Invoke(ctx, "/pos.PosReturnService/ReadAllPosReturns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PosReturnServiceServer is the server API for PosReturnService service.
// All implementations must embed UnimplementedPosReturnServiceServer
// for forward compatibility
type PosReturnServiceServer interface {
	CreatePosReturn(context.Context, *CreatePosReturnRequest) (*CreatePosReturnResponse, error)
	ReadPosReturn(context.Context, *ReadPosReturnRequest) (*ReadPosReturnResponse, error)
	UpdatePosReturn(context.Context, *UpdatePosReturnRequest) (*UpdatePosReturnResponse, error)
	DeletePosReturn(context.Context, *DeletePosReturnRequest) (*DeletePosReturnResponse, error)
	ReadAllPosReturns(context.Context, *ReadAllPosReturnsRequest) (*ReadAllPosReturnsResponse, error)
	mustEmbedUnimplementedPosReturnServiceServer()
}

// UnimplementedPosReturnServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPosReturnServiceServer struct {
}

func (UnimplementedPosReturnServiceServer) CreatePosReturn(context.Context, *CreatePosReturnRequest) (*CreatePosReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePosReturn not implemented")
}
func (UnimplementedPosReturnServiceServer) ReadPosReturn(context.Context, *ReadPosReturnRequest) (*ReadPosReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadPosReturn not implemented")
}
func (UnimplementedPosReturnServiceServer) UpdatePosReturn(context.Context, *UpdatePosReturnRequest) (*UpdatePosReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePosReturn not implemented")
}
func (UnimplementedPosReturnServiceServer) DeletePosReturn(context.Context, *DeletePosReturnRequest) (*DeletePosReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePosReturn not implemented")
}
func (UnimplementedPosReturnServiceServer) ReadAllPosReturns(context.Context, *ReadAllPosReturnsRequest) (*ReadAllPosReturnsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAllPosReturns not implemented")
}
func (UnimplementedPosReturnServiceServer) mustEmbedUnimplementedPosReturnServiceServer() {}

// UnsafePosReturnServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PosReturnServiceServer will
// result in compilation errors.
type UnsafePosReturnServiceServer interface {
	mustEmbedUnimplementedPosReturnServiceServer()
}

func RegisterPosReturnServiceServer(s grpc.ServiceRegistrar, srv PosReturnServiceServer) {
	s.RegisterService(&PosReturnService_ServiceDesc, srv)
}

func _PosReturnService_CreatePosReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePosReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosReturnServiceServer).CreatePosReturn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosReturnService/CreatePosReturn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosReturnServiceServer).CreatePosReturn(ctx, req.(*CreatePosReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosReturnService_ReadPosReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadPosReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosReturnServiceServer).ReadPosReturn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosReturnService/ReadPosReturn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosReturnServiceServer).ReadPosReturn(ctx, req.(*ReadPosReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosReturnService_UpdatePosReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePosReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosReturnServiceServer).UpdatePosReturn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosReturnService/UpdatePosReturn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosReturnServiceServer).UpdatePosReturn(ctx, req.(*UpdatePosReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosReturnService_DeletePosReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePosReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosReturnServiceServer).DeletePosReturn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosReturnService/DeletePosReturn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosReturnServiceServer).DeletePosReturn(ctx, req.(*DeletePosReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosReturnService_ReadAllPosReturns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAllPosReturnsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosReturnServiceServer).ReadAllPosReturns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosReturnService/ReadAllPosReturns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosReturnServiceServer).ReadAllPosReturns(ctx, req.(*ReadAllPosReturnsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PosReturnService_ServiceDesc is the grpc.ServiceDesc for PosReturnService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PosReturnService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pos.PosReturnService",
	HandlerType: (*PosReturnServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePosReturn",
			Handler:    _PosReturnService_CreatePosReturn_Handler,
		},
		{
			MethodName: "ReadPosReturn",
			Handler:    _PosReturnService_ReadPosReturn_Handler,
		},
		{
			MethodName: "UpdatePosReturn",
			Handler:    _PosReturnService_UpdatePosReturn_Handler,
		},
		{
			MethodName: "DeletePosReturn",
			Handler:    _PosReturnService_DeletePosReturn_Handler,
		},
		{
			MethodName: "ReadAllPosReturns",
			Handler:    _PosReturnService_ReadAllPosReturns_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "return.proto",
}
