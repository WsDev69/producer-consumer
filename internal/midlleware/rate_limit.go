package midlleware

import (
	"context"

	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(rl ratelimit.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		rl.Take()
		return handler(ctx, req)

		// it's not clear what to do when limit is reached, will leave it here
		// return nil, status.Error(codes.ResourceExhausted, "reached hit per limit")
	}
}
