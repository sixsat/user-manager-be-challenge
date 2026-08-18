package infra

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func loggerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	res, err := handler(ctx, req)
	slog.Info(
		"[grpc] request",
		slog.String("grpc_method", info.FullMethod),
		slog.Duration("latency", time.Since(start)),
		slog.String("grpc_code", status.Code(err).String()),
	)
	return res, err
}

func jwtAuthInterceptor(jwtSignKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(values) != 1 {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		scheme, tokenStr, found := strings.Cut(values[0], " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || tokenStr == "" {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		token, err := jwt.Parse(
			tokenStr,
			func(token *jwt.Token) (any, error) {
				return []byte(jwtSignKey), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithExpirationRequired(),
		)
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		return handler(ctx, req)
	}
}
