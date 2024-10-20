// Copyright 2023, NVIDIA CORPORATION & AFFILIATES. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//  * Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//  * Neither the name of NVIDIA CORPORATION nor the names of its
//    contributors may be used to endorse or promote products derived
//    from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS ``AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
// EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY
// OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: health.proto

package genconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	gen "grpc_test/gen"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// HealthName is the fully-qualified name of the Health service.
	HealthName = "Health"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// HealthCheckProcedure is the fully-qualified name of the Health's Check RPC.
	HealthCheckProcedure = "/Health/Check"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	healthServiceDescriptor     = gen.File_health_proto.Services().ByName("Health")
	healthCheckMethodDescriptor = healthServiceDescriptor.Methods().ByName("Check")
)

// HealthClient is a client for the Health service.
type HealthClient interface {
	// @@  .. cpp:var:: rpc Check(HealthCheckRequest) returns
	// @@       (HealthCheckResponse)
	// @@
	// @@     Get serving status of the inference server.
	// @@
	Check(context.Context, *connect.Request[gen.HealthCheckRequest]) (*connect.Response[gen.HealthCheckResponse], error)
}

// NewHealthClient constructs a client for the Health service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewHealthClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) HealthClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &healthClient{
		check: connect.NewClient[gen.HealthCheckRequest, gen.HealthCheckResponse](
			httpClient,
			baseURL+HealthCheckProcedure,
			connect.WithSchema(healthCheckMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// healthClient implements HealthClient.
type healthClient struct {
	check *connect.Client[gen.HealthCheckRequest, gen.HealthCheckResponse]
}

// Check calls Health.Check.
func (c *healthClient) Check(ctx context.Context, req *connect.Request[gen.HealthCheckRequest]) (*connect.Response[gen.HealthCheckResponse], error) {
	return c.check.CallUnary(ctx, req)
}

// HealthHandler is an implementation of the Health service.
type HealthHandler interface {
	// @@  .. cpp:var:: rpc Check(HealthCheckRequest) returns
	// @@       (HealthCheckResponse)
	// @@
	// @@     Get serving status of the inference server.
	// @@
	Check(context.Context, *connect.Request[gen.HealthCheckRequest]) (*connect.Response[gen.HealthCheckResponse], error)
}

// NewHealthHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewHealthHandler(svc HealthHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	healthCheckHandler := connect.NewUnaryHandler(
		HealthCheckProcedure,
		svc.Check,
		connect.WithSchema(healthCheckMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/Health/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case HealthCheckProcedure:
			healthCheckHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedHealthHandler returns CodeUnimplemented from all methods.
type UnimplementedHealthHandler struct{}

func (UnimplementedHealthHandler) Check(context.Context, *connect.Request[gen.HealthCheckRequest]) (*connect.Response[gen.HealthCheckResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("Health.Check is not implemented"))
}
