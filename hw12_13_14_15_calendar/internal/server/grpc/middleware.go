package internalgrpc

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func logginInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	addr := "unknown"

	start := time.Now()
	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	resp, err := handler(ctx, req)

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"loggingInterceptor",
		slog.String("addr", addr),
		slog.String("method", info.FullMethod),
		slog.Any("err", err),
		slog.String("duration", time.Since(start).String()),
	)

	return resp, err
}
