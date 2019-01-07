// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/apiserver/apiserver.proto

/*
Package apiserver is a generated protocol buffer package.

It is generated from these files:
	pkg/apiserver/apiserver.proto

It has these top-level messages:
	WorkflowList
	AddTaskRequest
	InvocationListQuery
	WorkflowInvocationList
	InvocationEventsResponse
	Health
*/
package apiserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import fission_workflows_types1 "github.com/fission/fission-workflows/pkg/types"
import fission_workflows_version "github.com/fission/fission-workflows/pkg/version"
import fission_workflows_eventstore "github.com/fission/fission-workflows/pkg/fes"
import google_protobuf2 "github.com/golang/protobuf/ptypes/empty"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WorkflowList struct {
	Workflows []string `protobuf:"bytes,1,rep,name=workflows" json:"workflows,omitempty"`
}

func (m *WorkflowList) Reset()                    { *m = WorkflowList{} }
func (m *WorkflowList) String() string            { return proto.CompactTextString(m) }
func (*WorkflowList) ProtoMessage()               {}
func (*WorkflowList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *WorkflowList) GetWorkflows() []string {
	if m != nil {
		return m.Workflows
	}
	return nil
}

type AddTaskRequest struct {
	InvocationID string                         `protobuf:"bytes,1,opt,name=invocationID" json:"invocationID,omitempty"`
	Task         *fission_workflows_types1.Task `protobuf:"bytes,2,opt,name=task" json:"task,omitempty"`
}

func (m *AddTaskRequest) Reset()                    { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()               {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AddTaskRequest) GetInvocationID() string {
	if m != nil {
		return m.InvocationID
	}
	return ""
}

func (m *AddTaskRequest) GetTask() *fission_workflows_types1.Task {
	if m != nil {
		return m.Task
	}
	return nil
}

type InvocationListQuery struct {
	Workflows []string `protobuf:"bytes,1,rep,name=workflows" json:"workflows,omitempty"`
}

func (m *InvocationListQuery) Reset()                    { *m = InvocationListQuery{} }
func (m *InvocationListQuery) String() string            { return proto.CompactTextString(m) }
func (*InvocationListQuery) ProtoMessage()               {}
func (*InvocationListQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *InvocationListQuery) GetWorkflows() []string {
	if m != nil {
		return m.Workflows
	}
	return nil
}

type WorkflowInvocationList struct {
	Invocations []string `protobuf:"bytes,1,rep,name=invocations" json:"invocations,omitempty"`
}

func (m *WorkflowInvocationList) Reset()                    { *m = WorkflowInvocationList{} }
func (m *WorkflowInvocationList) String() string            { return proto.CompactTextString(m) }
func (*WorkflowInvocationList) ProtoMessage()               {}
func (*WorkflowInvocationList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *WorkflowInvocationList) GetInvocations() []string {
	if m != nil {
		return m.Invocations
	}
	return nil
}

type InvocationEventsResponse struct {
	Metadata *fission_workflows_types1.ObjectMetadata `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	Events   []*fission_workflows_eventstore.Event    `protobuf:"bytes,2,rep,name=events" json:"events,omitempty"`
}

func (m *InvocationEventsResponse) Reset()                    { *m = InvocationEventsResponse{} }
func (m *InvocationEventsResponse) String() string            { return proto.CompactTextString(m) }
func (*InvocationEventsResponse) ProtoMessage()               {}
func (*InvocationEventsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *InvocationEventsResponse) GetMetadata() *fission_workflows_types1.ObjectMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *InvocationEventsResponse) GetEvents() []*fission_workflows_eventstore.Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type Health struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *Health) Reset()                    { *m = Health{} }
func (m *Health) String() string            { return proto.CompactTextString(m) }
func (*Health) ProtoMessage()               {}
func (*Health) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Health) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*WorkflowList)(nil), "fission.workflows.apiserver.WorkflowList")
	proto.RegisterType((*AddTaskRequest)(nil), "fission.workflows.apiserver.AddTaskRequest")
	proto.RegisterType((*InvocationListQuery)(nil), "fission.workflows.apiserver.InvocationListQuery")
	proto.RegisterType((*WorkflowInvocationList)(nil), "fission.workflows.apiserver.WorkflowInvocationList")
	proto.RegisterType((*InvocationEventsResponse)(nil), "fission.workflows.apiserver.InvocationEventsResponse")
	proto.RegisterType((*Health)(nil), "fission.workflows.apiserver.Health")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for WorkflowAPI service

type WorkflowAPIClient interface {
	Create(ctx context.Context, in *fission_workflows_types1.WorkflowSpec, opts ...grpc.CallOption) (*fission_workflows_types1.ObjectMetadata, error)
	List(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*WorkflowList, error)
	Get(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*fission_workflows_types1.Workflow, error)
	Delete(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	Validate(ctx context.Context, in *fission_workflows_types1.WorkflowSpec, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}

type workflowAPIClient struct {
	cc *grpc.ClientConn
}

func NewWorkflowAPIClient(cc *grpc.ClientConn) WorkflowAPIClient {
	return &workflowAPIClient{cc}
}

func (c *workflowAPIClient) Create(ctx context.Context, in *fission_workflows_types1.WorkflowSpec, opts ...grpc.CallOption) (*fission_workflows_types1.ObjectMetadata, error) {
	out := new(fission_workflows_types1.ObjectMetadata)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowAPI/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowAPIClient) List(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*WorkflowList, error) {
	out := new(WorkflowList)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowAPI/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowAPIClient) Get(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*fission_workflows_types1.Workflow, error) {
	out := new(fission_workflows_types1.Workflow)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowAPI/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowAPIClient) Delete(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowAPI/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowAPIClient) Validate(ctx context.Context, in *fission_workflows_types1.WorkflowSpec, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowAPI/Validate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for WorkflowAPI service

type WorkflowAPIServer interface {
	Create(context.Context, *fission_workflows_types1.WorkflowSpec) (*fission_workflows_types1.ObjectMetadata, error)
	List(context.Context, *google_protobuf2.Empty) (*WorkflowList, error)
	Get(context.Context, *fission_workflows_types1.ObjectMetadata) (*fission_workflows_types1.Workflow, error)
	Delete(context.Context, *fission_workflows_types1.ObjectMetadata) (*google_protobuf2.Empty, error)
	Validate(context.Context, *fission_workflows_types1.WorkflowSpec) (*google_protobuf2.Empty, error)
}

func RegisterWorkflowAPIServer(s *grpc.Server, srv WorkflowAPIServer) {
	s.RegisterService(&_WorkflowAPI_serviceDesc, srv)
}

func _WorkflowAPI_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.WorkflowSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowAPIServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowAPI/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowAPIServer).Create(ctx, req.(*fission_workflows_types1.WorkflowSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowAPI_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf2.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowAPIServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowAPI/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowAPIServer).List(ctx, req.(*google_protobuf2.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowAPI_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.ObjectMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowAPIServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowAPI/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowAPIServer).Get(ctx, req.(*fission_workflows_types1.ObjectMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowAPI_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.ObjectMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowAPIServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowAPI/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowAPIServer).Delete(ctx, req.(*fission_workflows_types1.ObjectMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowAPI_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.WorkflowSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowAPIServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowAPI/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowAPIServer).Validate(ctx, req.(*fission_workflows_types1.WorkflowSpec))
	}
	return interceptor(ctx, in, info, handler)
}

var _WorkflowAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fission.workflows.apiserver.WorkflowAPI",
	HandlerType: (*WorkflowAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _WorkflowAPI_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _WorkflowAPI_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _WorkflowAPI_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _WorkflowAPI_Delete_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _WorkflowAPI_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apiserver/apiserver.proto",
}

// Client API for WorkflowInvocationAPI service

type WorkflowInvocationAPIClient interface {
	// Create a new workflow invocation
	//
	// In case the invocation specification is missing fields or contains invalid fields, a HTTP 400 is returned.
	Invoke(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*fission_workflows_types1.ObjectMetadata, error)
	InvokeSync(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*fission_workflows_types1.WorkflowInvocation, error)
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// Cancel a workflow invocation
	//
	// This action is irreverisble. A canceled invocation cannot be resumed or restarted.
	// In case that an invocation already is canceled, has failed or has completed, nothing happens.
	// In case that an invocation does not exist a HTTP 404 error status is returned.
	Cancel(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	List(ctx context.Context, in *InvocationListQuery, opts ...grpc.CallOption) (*WorkflowInvocationList, error)
	// Get the specification and status of a workflow invocation
	//
	// Get returns three different aspects of the workflow invocation, namely the spec (specification), status and logs.
	// To lighten the request load, consider using a more specific request.
	Get(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*fission_workflows_types1.WorkflowInvocation, error)
	Events(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*InvocationEventsResponse, error)
	Validate(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}

type workflowInvocationAPIClient struct {
	cc *grpc.ClientConn
}

func NewWorkflowInvocationAPIClient(cc *grpc.ClientConn) WorkflowInvocationAPIClient {
	return &workflowInvocationAPIClient{cc}
}

func (c *workflowInvocationAPIClient) Invoke(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*fission_workflows_types1.ObjectMetadata, error) {
	out := new(fission_workflows_types1.ObjectMetadata)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/Invoke", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) InvokeSync(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*fission_workflows_types1.WorkflowInvocation, error) {
	out := new(fission_workflows_types1.WorkflowInvocation)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/InvokeSync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/AddTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) Cancel(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/Cancel", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) List(ctx context.Context, in *InvocationListQuery, opts ...grpc.CallOption) (*WorkflowInvocationList, error) {
	out := new(WorkflowInvocationList)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) Get(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*fission_workflows_types1.WorkflowInvocation, error) {
	out := new(fission_workflows_types1.WorkflowInvocation)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) Events(ctx context.Context, in *fission_workflows_types1.ObjectMetadata, opts ...grpc.CallOption) (*InvocationEventsResponse, error) {
	out := new(InvocationEventsResponse)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/Events", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowInvocationAPIClient) Validate(ctx context.Context, in *fission_workflows_types1.WorkflowInvocationSpec, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.WorkflowInvocationAPI/Validate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for WorkflowInvocationAPI service

type WorkflowInvocationAPIServer interface {
	// Create a new workflow invocation
	//
	// In case the invocation specification is missing fields or contains invalid fields, a HTTP 400 is returned.
	Invoke(context.Context, *fission_workflows_types1.WorkflowInvocationSpec) (*fission_workflows_types1.ObjectMetadata, error)
	InvokeSync(context.Context, *fission_workflows_types1.WorkflowInvocationSpec) (*fission_workflows_types1.WorkflowInvocation, error)
	AddTask(context.Context, *AddTaskRequest) (*google_protobuf2.Empty, error)
	// Cancel a workflow invocation
	//
	// This action is irreverisble. A canceled invocation cannot be resumed or restarted.
	// In case that an invocation already is canceled, has failed or has completed, nothing happens.
	// In case that an invocation does not exist a HTTP 404 error status is returned.
	Cancel(context.Context, *fission_workflows_types1.ObjectMetadata) (*google_protobuf2.Empty, error)
	List(context.Context, *InvocationListQuery) (*WorkflowInvocationList, error)
	// Get the specification and status of a workflow invocation
	//
	// Get returns three different aspects of the workflow invocation, namely the spec (specification), status and logs.
	// To lighten the request load, consider using a more specific request.
	Get(context.Context, *fission_workflows_types1.ObjectMetadata) (*fission_workflows_types1.WorkflowInvocation, error)
	Events(context.Context, *fission_workflows_types1.ObjectMetadata) (*InvocationEventsResponse, error)
	Validate(context.Context, *fission_workflows_types1.WorkflowInvocationSpec) (*google_protobuf2.Empty, error)
}

func RegisterWorkflowInvocationAPIServer(s *grpc.Server, srv WorkflowInvocationAPIServer) {
	s.RegisterService(&_WorkflowInvocationAPI_serviceDesc, srv)
}

func _WorkflowInvocationAPI_Invoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.WorkflowInvocationSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).Invoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/Invoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).Invoke(ctx, req.(*fission_workflows_types1.WorkflowInvocationSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_InvokeSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.WorkflowInvocationSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).InvokeSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/InvokeSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).InvokeSync(ctx, req.(*fission_workflows_types1.WorkflowInvocationSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/AddTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_Cancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.ObjectMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).Cancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/Cancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).Cancel(ctx, req.(*fission_workflows_types1.ObjectMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvocationListQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).List(ctx, req.(*InvocationListQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.ObjectMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).Get(ctx, req.(*fission_workflows_types1.ObjectMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_Events_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.ObjectMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).Events(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/Events",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).Events(ctx, req.(*fission_workflows_types1.ObjectMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowInvocationAPI_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(fission_workflows_types1.WorkflowInvocationSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowInvocationAPIServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.WorkflowInvocationAPI/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowInvocationAPIServer).Validate(ctx, req.(*fission_workflows_types1.WorkflowInvocationSpec))
	}
	return interceptor(ctx, in, info, handler)
}

var _WorkflowInvocationAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fission.workflows.apiserver.WorkflowInvocationAPI",
	HandlerType: (*WorkflowInvocationAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Invoke",
			Handler:    _WorkflowInvocationAPI_Invoke_Handler,
		},
		{
			MethodName: "InvokeSync",
			Handler:    _WorkflowInvocationAPI_InvokeSync_Handler,
		},
		{
			MethodName: "AddTask",
			Handler:    _WorkflowInvocationAPI_AddTask_Handler,
		},
		{
			MethodName: "Cancel",
			Handler:    _WorkflowInvocationAPI_Cancel_Handler,
		},
		{
			MethodName: "List",
			Handler:    _WorkflowInvocationAPI_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _WorkflowInvocationAPI_Get_Handler,
		},
		{
			MethodName: "Events",
			Handler:    _WorkflowInvocationAPI_Events_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _WorkflowInvocationAPI_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apiserver/apiserver.proto",
}

// Client API for AdminAPI service

type AdminAPIClient interface {
	Status(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*Health, error)
	Version(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*fission_workflows_version.Info, error)
}

type adminAPIClient struct {
	cc *grpc.ClientConn
}

func NewAdminAPIClient(cc *grpc.ClientConn) AdminAPIClient {
	return &adminAPIClient{cc}
}

func (c *adminAPIClient) Status(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*Health, error) {
	out := new(Health)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.AdminAPI/Status", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminAPIClient) Version(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*fission_workflows_version.Info, error) {
	out := new(fission_workflows_version.Info)
	err := grpc.Invoke(ctx, "/fission.workflows.apiserver.AdminAPI/Version", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AdminAPI service

type AdminAPIServer interface {
	Status(context.Context, *google_protobuf2.Empty) (*Health, error)
	Version(context.Context, *google_protobuf2.Empty) (*fission_workflows_version.Info, error)
}

func RegisterAdminAPIServer(s *grpc.Server, srv AdminAPIServer) {
	s.RegisterService(&_AdminAPI_serviceDesc, srv)
}

func _AdminAPI_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf2.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminAPIServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.AdminAPI/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminAPIServer).Status(ctx, req.(*google_protobuf2.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminAPI_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf2.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminAPIServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fission.workflows.apiserver.AdminAPI/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminAPIServer).Version(ctx, req.(*google_protobuf2.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdminAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fission.workflows.apiserver.AdminAPI",
	HandlerType: (*AdminAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _AdminAPI_Status_Handler,
		},
		{
			MethodName: "Version",
			Handler:    _AdminAPI_Version_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apiserver/apiserver.proto",
}

func init() { proto.RegisterFile("pkg/apiserver/apiserver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 813 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x51, 0x4f, 0xdb, 0x56,
	0x14, 0x96, 0x01, 0x99, 0xe4, 0x9a, 0xa1, 0xec, 0xc0, 0x42, 0x16, 0x40, 0x64, 0x17, 0x4d, 0x0b,
	0x61, 0xb3, 0xb7, 0xa0, 0xed, 0x21, 0x93, 0x26, 0x31, 0x40, 0x5b, 0xa4, 0x56, 0xb4, 0x01, 0x81,
	0x84, 0xfa, 0xe2, 0x38, 0x37, 0x89, 0x9b, 0xc4, 0x76, 0x7d, 0x6f, 0x82, 0x02, 0xe2, 0x85, 0xb7,
	0x3e, 0xf3, 0xd8, 0x4a, 0xfd, 0x17, 0x95, 0xfa, 0x3b, 0xfa, 0x17, 0xfa, 0x43, 0x2a, 0x5f, 0x5f,
	0xc7, 0x0e, 0xc1, 0x89, 0xa3, 0x3e, 0x80, 0xe3, 0xe3, 0xf3, 0x9d, 0xef, 0x3b, 0xc7, 0xe7, 0x9c,
	0x6b, 0xb4, 0xed, 0x74, 0x5a, 0x9a, 0xee, 0x98, 0x94, 0xb8, 0x03, 0xe2, 0x86, 0xbf, 0x54, 0xc7,
	0xb5, 0x99, 0x0d, 0x9b, 0x4d, 0x93, 0x52, 0xd3, 0xb6, 0xd4, 0x6b, 0xdb, 0xed, 0x34, 0xbb, 0xf6,
	0x35, 0x55, 0x47, 0x2e, 0xf9, 0x4a, 0xcb, 0x64, 0xed, 0x7e, 0x5d, 0x35, 0xec, 0x9e, 0x26, 0xfc,
	0x82, 0xeb, 0x6f, 0x23, 0x7f, 0xcd, 0x23, 0x60, 0x43, 0x87, 0x50, 0xff, 0xbf, 0x1f, 0x38, 0xff,
	0x4f, 0x62, 0xec, 0x80, 0xb8, 0xfc, 0xa9, 0xb8, 0x0a, 0xfc, 0x5f, 0x89, 0xf1, 0x4d, 0x42, 0xbd,
	0x3f, 0x81, 0xdb, 0x6c, 0xd9, 0x76, 0xab, 0x4b, 0x34, 0x7e, 0x57, 0xef, 0x37, 0x35, 0xd2, 0x73,
	0xd8, 0x50, 0x3c, 0xdc, 0x12, 0x0f, 0x75, 0xc7, 0xd4, 0x74, 0xcb, 0xb2, 0x99, 0xce, 0x4c, 0xdb,
	0x12, 0x50, 0xfc, 0x2b, 0x5a, 0xb9, 0x14, 0x91, 0x9f, 0x99, 0x94, 0xc1, 0x16, 0x4a, 0x8f, 0x98,
	0x72, 0x52, 0x61, 0xb1, 0x98, 0xae, 0x85, 0x06, 0xdc, 0x42, 0xab, 0x87, 0x8d, 0xc6, 0xb9, 0x4e,
	0x3b, 0x35, 0xf2, 0xa6, 0x4f, 0x28, 0x03, 0x8c, 0x56, 0x4c, 0x6b, 0x60, 0x1b, 0x3c, 0x68, 0xf5,
	0x38, 0x27, 0x15, 0xa4, 0x62, 0xba, 0x36, 0x66, 0x83, 0x3f, 0xd0, 0x12, 0xd3, 0x69, 0x27, 0xb7,
	0x50, 0x90, 0x8a, 0x4a, 0x79, 0x5b, 0x9d, 0x2c, 0xbf, 0x5f, 0x44, 0x1e, 0x97, 0xbb, 0xe2, 0x03,
	0xb4, 0x56, 0x1d, 0x85, 0xf0, 0x84, 0xbd, 0xec, 0x13, 0x77, 0x38, 0x43, 0x5d, 0x05, 0x65, 0x83,
	0x5c, 0xc6, 0xc1, 0x50, 0x40, 0x4a, 0xa8, 0x28, 0x40, 0x46, 0x4d, 0xf8, 0xbd, 0x84, 0x72, 0x21,
	0xe8, 0x64, 0x40, 0x2c, 0x46, 0x6b, 0x84, 0x3a, 0xb6, 0x45, 0x09, 0x1c, 0xa1, 0x54, 0x8f, 0x30,
	0xbd, 0xa1, 0x33, 0x9d, 0x27, 0xa8, 0x94, 0x7f, 0x89, 0x4d, 0xe2, 0xb4, 0xfe, 0x9a, 0x18, 0xec,
	0xb9, 0x70, 0xaf, 0x8d, 0x80, 0xf0, 0x37, 0x92, 0x09, 0x0f, 0x9b, 0x5b, 0x28, 0x2c, 0x16, 0x95,
	0xf2, 0xee, 0x13, 0x21, 0x7c, 0x07, 0x66, 0xbb, 0x44, 0xe5, 0x12, 0x6a, 0x02, 0x82, 0x0b, 0x48,
	0xfe, 0x9f, 0xe8, 0x5d, 0xd6, 0x86, 0x2c, 0x92, 0x29, 0xd3, 0x59, 0x9f, 0x8a, 0x52, 0x8b, 0xbb,
	0xf2, 0xc3, 0x12, 0x52, 0x82, 0xec, 0x0f, 0x5f, 0x54, 0xc1, 0x42, 0xf2, 0x91, 0x4b, 0x74, 0x46,
	0xe0, 0xe7, 0x58, 0xad, 0x81, 0xff, 0x99, 0x43, 0x8c, 0x7c, 0xd2, 0x94, 0xf0, 0xfa, 0xfd, 0xe7,
	0x2f, 0x0f, 0x0b, 0xab, 0x38, 0xad, 0x05, 0x8e, 0x15, 0xa9, 0x04, 0xaf, 0xd0, 0x12, 0x2f, 0x75,
	0x56, 0xf5, 0xfb, 0x4d, 0x0d, 0x9a, 0x51, 0x3d, 0xf1, 0x9a, 0x31, 0xbf, 0xa7, 0x4e, 0x99, 0x3a,
	0x35, 0xda, 0x83, 0xf8, 0x7b, 0x4e, 0xa0, 0x40, 0x48, 0x00, 0x26, 0x5a, 0xfc, 0x8f, 0x30, 0x48,
	0xaa, 0x31, 0xff, 0xd3, 0xcc, 0x9c, 0x71, 0x96, 0xb3, 0x64, 0x60, 0x75, 0xc4, 0xa2, 0xdd, 0x9a,
	0x8d, 0x3b, 0xd0, 0x91, 0x7c, 0x4c, 0xba, 0x84, 0x91, 0xe4, 0x6c, 0x31, 0x39, 0x07, 0x14, 0xa5,
	0xc7, 0x14, 0x6d, 0x94, 0xba, 0xd0, 0xbb, 0x66, 0x63, 0x8e, 0xb7, 0x13, 0x47, 0xb1, 0xcd, 0x29,
	0x36, 0x30, 0x84, 0x14, 0x03, 0x11, 0xba, 0x22, 0x95, 0xca, 0x1f, 0x53, 0xe8, 0x87, 0xc9, 0x99,
	0xf0, 0xfa, 0xe3, 0x06, 0xc9, 0x9e, 0xa1, 0x43, 0x40, 0x9b, 0xa9, 0x20, 0x44, 0xce, 0xd7, 0x29,
	0x22, 0x7f, 0xac, 0x68, 0xe1, 0xa8, 0x79, 0xbd, 0xf2, 0x4e, 0x42, 0xc8, 0x27, 0x3f, 0x1b, 0x5a,
	0xc6, 0xfc, 0x02, 0xf6, 0xe7, 0x00, 0x60, 0x8d, 0x8b, 0xd8, 0xc3, 0x99, 0x88, 0x08, 0x8d, 0x0e,
	0x2d, 0xa3, 0x22, 0x95, 0xae, 0x00, 0x26, 0xcc, 0xf0, 0x41, 0x42, 0xcb, 0x62, 0xcb, 0xc1, 0xfe,
	0xd4, 0xae, 0x1d, 0xdf, 0x85, 0xb1, 0xef, 0xe8, 0x94, 0x2b, 0xa8, 0xe2, 0x42, 0x94, 0xea, 0x36,
	0xba, 0x22, 0xef, 0x34, 0x6f, 0xeb, 0x51, 0x4f, 0x11, 0xce, 0xcf, 0x74, 0x03, 0x03, 0xc9, 0x47,
	0xba, 0x65, 0x90, 0xee, 0xb7, 0xb7, 0x68, 0x8e, 0x6b, 0x83, 0x52, 0x66, 0x9c, 0xb4, 0x71, 0x07,
	0xf7, 0x92, 0x98, 0xe8, 0xdf, 0xa7, 0xd6, 0xe0, 0x89, 0x35, 0x9d, 0x3f, 0x48, 0x34, 0xeb, 0xe3,
	0x48, 0xbc, 0xc6, 0x95, 0x7c, 0x07, 0xd1, 0x66, 0x81, 0xfe, 0x9c, 0x73, 0x3f, 0x57, 0x67, 0x88,
	0xdc, 0x61, 0x32, 0xf7, 0xb7, 0x12, 0x92, 0xfd, 0x33, 0x20, 0x39, 0xf5, 0x9f, 0x09, 0xcb, 0x34,
	0x7e, 0xb6, 0xe0, 0x1d, 0x2e, 0xe2, 0x47, 0xd8, 0x78, 0x2c, 0x42, 0xf3, 0x57, 0x3f, 0xb0, 0xc8,
	0xb2, 0x98, 0x7b, 0x52, 0xe2, 0x5e, 0xbb, 0x60, 0xc5, 0xeb, 0x51, 0xd6, 0xe8, 0xe2, 0xf8, 0x24,
	0xa1, 0xd4, 0x61, 0xa3, 0x67, 0xf2, 0x5d, 0x71, 0x89, 0xe4, 0x33, 0x7e, 0xca, 0xc4, 0x6e, 0xf7,
	0xdd, 0xa9, 0xc9, 0xfb, 0x47, 0x17, 0xce, 0x70, 0x52, 0x04, 0x29, 0xad, 0xcd, 0x0d, 0x37, 0x70,
	0x8e, 0x96, 0x2f, 0xfc, 0x2f, 0xa0, 0xd8, 0xc8, 0x3b, 0x4f, 0x44, 0x0e, 0xbe, 0x9a, 0xaa, 0x56,
	0xd3, 0x8e, 0x44, 0x15, 0xe6, 0x7f, 0x95, 0xab, 0xf4, 0x88, 0xbb, 0x2e, 0xf3, 0x78, 0x07, 0x5f,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x39, 0xda, 0x8d, 0x48, 0x14, 0x0a, 0x00, 0x00,
}
