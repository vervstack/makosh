package interceptors

import (
	"context"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthHeader - header containing auth token to talk to service
// in order to perform REST call should pre append runtime.MetadataHeaderPrefix (e.g - "Grpc-Metadata-")
const (
	AuthHeader = "Makosh-Auth"

	NoAuthErrMessage      = "no auth header"
	InvalidAuthErrMessage = "invalid auth header"
)

func GrpcAuthInterceptor(secret string) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.FailedPrecondition, "error unmarshalling metadata from context")
			}

			auth := md.Get(AuthHeader)
			if len(auth) == 0 {
				return nil, status.Error(codes.PermissionDenied, NoAuthErrMessage)
			}

			if !slices.Contains(auth, secret) {
				return nil, status.Error(codes.PermissionDenied, InvalidAuthErrMessage)
			}

			return handler(ctx, req)
		})
}
