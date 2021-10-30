package grpc_error

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/nhatminhk63j/uetvoting/config"
	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

// UnaryServerInterceptor returns a new unary server interceptor that wraps output error.
//
// Output error will be converted to GRPC error before sending to clients.

func UnaryServerInterceptor(appMode config.AppMode, internalServerErr error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if err != nil {
			go sentry.CaptureException(err)
			return nil, grpcErrorWrapper{appMode: appMode, internalServerErr: internalServerErr}.GRPCError(err)
		}
		return res, nil
	}
}

// grpcErrorWrapper is wrapper that convert app level error to GRPC error
type grpcErrorWrapper struct {
	appMode           config.AppMode
	internalServerErr error
}

// GRPCError converts original error to GRPC error which will then be converted to HTTP error by grpc-gateway.
// Error may be wrapped, so must unwrap it to retrieve original error.
func (w grpcErrorWrapper) GRPCError(err error) error {
	wrappedErr := unwrapErr(err)
	if wrappedErr == context.Canceled || wrappedErr == context.DeadlineExceeded {
		return status.FromContextError(wrappedErr).Err()
	}
	stt, ok := status.FromError(wrappedErr)

	if de, containDetails := wrappedErr.(interface {
		Details() []proto.Message
	}); containDetails {
		if s, err := stt.WithDetails(de.Details()...); err == nil {
			stt = s
		}
	}

	if ok {
		return stt.Err()
	}

	// In development mode, return raw error message.
	if w.appMode == config.DevelopmentMode {
		logger.WithFields(logger.Fields{"error": err}).Warnf("getting error...")
		return status.Error(stt.Code(), err.Error())
	}

	logger.WithFields(logger.Fields{"error": err}).Error("unexpected error...")
	return w.internalServerErr
}

func unwrapErr(err error) error {
	for err != nil {
		wrappedError := errors.Unwrap(err)
		if wrappedError == nil {
			break
		}
		err = wrappedError
	}

	return err
}
