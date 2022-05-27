package services

import (
	"context"
	"fmt"

	"github.com/wpcodevo/golang-mongodb/config"
	"github.com/wpcodevo/golang-mongodb/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	allowedRoles map[string][]string
	config       *config.Config
	userService  UserService
}

func NewAuthInterceptor(allowedRoles map[string][]string, config *config.Config, userService UserService) *AuthInterceptor {
	return &AuthInterceptor{allowedRoles, config, userService}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		fmt.Println("Unary interceptor:", info.FullMethod)
		err := interceptor.authorize(ctx, info.FullMethod)

		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	allowedRoles, ok := interceptor.allowedRoles[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Errorf(codes.Unauthenticated, "You are not logged in")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "access token not provided")

	}

	accessToken := values[0]

	userId, err := utils.ValidateToken(accessToken, interceptor.config.AccessTokenPublicKey)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	user, err := interceptor.userService.FindUserById(fmt.Sprint(userId))
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}

	for _, role := range allowedRoles {
		if role == user.Role {
			return nil
		}

	}
	return status.Errorf(codes.PermissionDenied, "You are not permitted to access this RPC")

}
