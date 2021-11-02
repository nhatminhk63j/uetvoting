package authorization

import (
	"context"
	"google.golang.org/grpc"

	"github.com/nhatminhk63j/uetvoting/server/middleware/authentication"
)

const (
	UserPermission  = 1
	AdminPermission = 2

	//	Prefix route
	eventServiceRoute string = "/event.v1.EventService/"
)

// accessibleRoles define permission for user request.
var accessibleRoles = map[string]int{
	eventServiceRoute + "UpsertEvent": UserPermission,
}

// UnaryServerInterceptor to check authorization.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := authorize(ctx, info)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Check authorization of role user by accessToken from header request.
func authorize(ctx context.Context, info *grpc.UnaryServerInfo) (err error) {
	method := info.FullMethod
	accessibleRole, ok := accessibleRoles[method] // fix here to get route and get permission for this
	if !ok {
		// everyone can access
		return
	}

	userInfo := authentication.GetUserInfoFromContext(ctx)
	if userInfo == nil {
		return NoPermissionError{}
	}

	if userInfo.CheckRole(accessibleRole) {
		return nil
	}

	return NoPermissionError{}
}
