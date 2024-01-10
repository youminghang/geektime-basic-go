package trace

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx/interceptors"
	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type OTELInterceptorBuilder struct {
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
	interceptors.Builder
	serviceName string
}

func NewOTELInterceptorBuilder(
	serviceName string,
	tracer trace.Tracer,
	propagator propagation.TextMapPropagator) *OTELInterceptorBuilder {
	return &OTELInterceptorBuilder{tracer: tracer, propagator: propagator}
}

func (b *OTELInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.GetTracerProvider().
			Tracer("gitee.com/geekbang/basic-go/webook/pkg/grpcx")
	}
	propagator := b.propagator
	if propagator == nil {
		propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		ctx, span := tracer.Start(ctx, info.FullMethod, trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(b.PeerName(ctx)),
			attribute.Key("net.peer.ip").String(b.PeerIP(ctx)),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return handler(ctx, req)
	}
}

func (b *OTELInterceptorBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.GetTracerProvider().
			Tracer("gitee.com/geekbang/basic-go/webook/pkg/grpcx")
	}
	propagator := b.propagator
	if propagator == nil {
		propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		ctx, span := tracer.Start(ctx, method, trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(method),
			semconv.NetPeerNameKey.String(b.serviceName),
		)
		// 把这个往后传递
		ctx = metadata.NewOutgoingContext(ctx, md)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
