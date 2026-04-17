package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pdaccess/commons/pkg/domain"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func authzInterceptor(
	isExemptMethod func(string) bool,
	validator func(context.Context, string) (*domain.PdaccessClaims, error),
) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if isExemptMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr, found := strings.CutPrefix(authHeader[0], "Bearer ")
		if !found {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
		}

		pdaccessClaims, err := validator(ctx, tokenStr)
		if err != nil {
			return nil, fmt.Errorf("token validation: %w", err)
		}

		clientRequest := &ClientRequest{
			User:      *pdaccessClaims,
			RequestId: ClientFromCtx(ctx).RequestId,
		}
		ctx = context.WithValue(ctx, ClientRequestType(ClientRequestStr), clientRequest)

		return handler(ctx, req)
	}
}

func loggingInterceptor(isExemptMethod func(string) bool) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		methodName := extractMethodName(info.FullMethod)
		tokenInfo := ClientFromCtx(ctx)

		log.Info().
			Str("method", methodName).
			Str("full_method", info.FullMethod).
			Str("user_id", tokenInfo.User.UserId).
			Str("user_role", tokenInfo.User.Role).
			Str("realm", tokenInfo.User.Realm).
			Str("auth_id", tokenInfo.User.AuthId).
			Time("start_time", startTime).
			Msg("gRPC request received")

		if isExemptMethod(info.FullMethod) {
			log.Debug().Interface("request", req).Msg("gRPC request data")
		}

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			log.Error().
				Str("method", methodName).
				Str("full_method", info.FullMethod).
				Str("user_id", tokenInfo.User.UserId).
				Str("error", err.Error()).
				Dur("duration", duration).
				Msg("gRPC request failed")
		} else {
			log.Info().
				Str("method", methodName).
				Str("full_method", info.FullMethod).
				Str("user_id", tokenInfo.User.UserId).
				Dur("duration", duration).
				Msg("gRPC request completed")
		}

		return resp, err
	}
}

func loggingStreamInterceptor(isExemptMethod func(string) bool) func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		methodName := extractMethodName(info.FullMethod)
		tokenInfo := ClientFromCtx(ss.Context())

		log.Info().
			Str("method", methodName).
			Str("full_method", info.FullMethod).
			Str("user_id", tokenInfo.User.UserId).
			Str("user_role", tokenInfo.User.Role).
			Str("realm", tokenInfo.User.Realm).
			Time("start_time", startTime).
			Msg("gRPC stream request received")

		err := handler(srv, ss)

		duration := time.Since(startTime)
		if err != nil {
			log.Error().
				Str("method", methodName).
				Str("full_method", info.FullMethod).
				Str("user_id", tokenInfo.User.UserId).
				Str("error", err.Error()).
				Dur("duration", duration).
				Msg("gRPC stream request failed")
		} else {
			log.Info().
				Str("method", methodName).
				Str("full_method", info.FullMethod).
				Str("user_id", tokenInfo.User.UserId).
				Dur("duration", duration).
				Msg("gRPC stream request completed")
		}

		return err
	}
}

func extractMethodName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullMethod
}
