// Copyright (c) 2022 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: section.proto

package section_v1

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

// SectionClient is the client API for Section service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SectionClient interface {
	// *
	// GetRotationItems: get current rotation items, this method will be called by rotation type is CUSTOM
	GetRotationItems(ctx context.Context, in *GetRotationItemsRequest, opts ...grpc.CallOption) (*GetRotationItemsResponse, error)
	// *
	// Backfill method trigger condition:
	// 1. Rotation type is FIXED_PERIOD
	// 2. Backfill type is CUSTOM
	// 3. User already owned any one of current rotation items.
	Backfill(ctx context.Context, in *BackfillRequest, opts ...grpc.CallOption) (*BackfillResponse, error)
}

type sectionClient struct {
	cc grpc.ClientConnInterface
}

func NewSectionClient(cc grpc.ClientConnInterface) SectionClient {
	return &sectionClient{cc}
}

func (c *sectionClient) GetRotationItems(ctx context.Context, in *GetRotationItemsRequest, opts ...grpc.CallOption) (*GetRotationItemsResponse, error) {
	out := new(GetRotationItemsResponse)
	err := c.cc.Invoke(ctx, "/accelbyte.platform.catalog.section.v1.Section/GetRotationItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sectionClient) Backfill(ctx context.Context, in *BackfillRequest, opts ...grpc.CallOption) (*BackfillResponse, error) {
	out := new(BackfillResponse)
	err := c.cc.Invoke(ctx, "/accelbyte.platform.catalog.section.v1.Section/Backfill", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SectionServer is the server API for Section service.
// All implementations must embed UnimplementedSectionServer
// for forward compatibility
type SectionServer interface {
	// *
	// GetRotationItems: get current rotation items, this method will be called by rotation type is CUSTOM
	GetRotationItems(context.Context, *GetRotationItemsRequest) (*GetRotationItemsResponse, error)
	// *
	// Backfill method trigger condition:
	// 1. Rotation type is FIXED_PERIOD
	// 2. Backfill type is CUSTOM
	// 3. User already owned any one of current rotation items.
	Backfill(context.Context, *BackfillRequest) (*BackfillResponse, error)
	mustEmbedUnimplementedSectionServer()
}

// UnimplementedSectionServer must be embedded to have forward compatible implementations.
type UnimplementedSectionServer struct {
}

func (UnimplementedSectionServer) GetRotationItems(context.Context, *GetRotationItemsRequest) (*GetRotationItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRotationItems not implemented")
}
func (UnimplementedSectionServer) Backfill(context.Context, *BackfillRequest) (*BackfillResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Backfill not implemented")
}
func (UnimplementedSectionServer) mustEmbedUnimplementedSectionServer() {}

// UnsafeSectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SectionServer will
// result in compilation errors.
type UnsafeSectionServer interface {
	mustEmbedUnimplementedSectionServer()
}

func RegisterSectionServer(s grpc.ServiceRegistrar, srv SectionServer) {
	s.RegisterService(&Section_ServiceDesc, srv)
}

func _Section_GetRotationItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRotationItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SectionServer).GetRotationItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accelbyte.platform.catalog.section.v1.Section/GetRotationItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SectionServer).GetRotationItems(ctx, req.(*GetRotationItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Section_Backfill_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackfillRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SectionServer).Backfill(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accelbyte.platform.catalog.section.v1.Section/Backfill",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SectionServer).Backfill(ctx, req.(*BackfillRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Section_ServiceDesc is the grpc.ServiceDesc for Section service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Section_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "accelbyte.platform.catalog.section.v1.Section",
	HandlerType: (*SectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRotationItems",
			Handler:    _Section_GetRotationItems_Handler,
		},
		{
			MethodName: "Backfill",
			Handler:    _Section_Backfill_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "section.proto",
}