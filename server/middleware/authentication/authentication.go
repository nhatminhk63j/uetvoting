package authentication

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/nhatminhk63j/uetvoting/pkg/jwt"
)

const (
	userInfoCtxKey         = "userInfoCtx"
	authorizationHeaderKey = "authorization"
)

// UserInfo contains user info.
type UserInfo struct {
	Id    int
	Email string
	Role  int
}

// CheckRole ...
func (u *UserInfo) CheckRole(role int) bool {
	return (u.Role & role) > 0
}

// UnaryServerInterceptor ...
func UnaryServerInterceptor(jwtManager *jwt.Resolver) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		userInfo, err := authenticate(ctx, jwtManager)
		if err == nil {
			ctx = context.WithValue(ctx, userInfoCtxKey, userInfo)
		} else {
			var errType *MissingMetadataError
			if errors.As(err, &errType) {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

// GetUserInfoFromContext ...
func GetUserInfoFromContext(ctx context.Context) *UserInfo {
	if v := ctx.Value(userInfoCtxKey); v != nil {
		if info, ok := v.(*UserInfo); ok {
			return info
		}
	}
	return nil
}

func authenticate(ctx context.Context, jwtManager *jwt.Resolver) (*UserInfo, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, MissingMetadataError{}
	}

	values := md[authorizationHeaderKey]
	if len(values) == 0 {
		return nil, MissingAccessTokenError{}
	}

	accessToken := values[0]
	userClaims, err := jwtManager.Verify(accessToken)
	if err != nil {
		return nil, InvalidAccessTokenError{}
	}

	return &UserInfo{
		Id:    userClaims.ID,
		Email: userClaims.Email,
		Role:  userClaims.Role,
	}, nil
}
