package interceptors

import (
	"connect-crm-backend/crm-connect-backend/internal/config"
	"connect-crm-backend/crm-connect-backend/internal/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

import (
	"context"
)

type authInterceptor struct {
	*config.MiddlewareConfig
}

func NewAuthInterceptor(middlewareConfig *config.MiddlewareConfig) *authInterceptor {
	return &authInterceptor{
		middlewareConfig,
	}
}

func (a *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logrus.Debugf("AuthInterceptor Unary Interceptor : validating bearer token... for method %v", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Meta Data is not provided")
		}

		values := md.Get("Authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Authentication Token Missing")
		}

		token := values[0]
		authSplits := strings.Split(token, " ")
		if len(authSplits) < 2 {
			return nil, status.Error(codes.Unauthenticated, "Authentication Token Missing")
		}

		authToken := strings.TrimSpace(authSplits[1])
		enabledTokens := a.EnabledTokens
		if util.Contains(enabledTokens, authToken) < 0 {
			return nil, status.Error(codes.Unauthenticated, "Authentication Token passed is invalid")
		}

		logrus.Debug("AuthInterceptor Unary Interceptor : Token Validation successful...")

		return handler(ctx, req)
	}
}
