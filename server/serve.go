package server

import (
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nhatminhk63j/uetvoting/config"
	"github.com/nhatminhk63j/uetvoting/pkg/jwt"
	"github.com/nhatminhk63j/uetvoting/pkg/logger"
	"github.com/nhatminhk63j/uetvoting/server/handler/auth"
	"github.com/nhatminhk63j/uetvoting/server/handler/event"
	"github.com/nhatminhk63j/uetvoting/server/handler/health"
	"github.com/nhatminhk63j/uetvoting/server/middleware/authentication"
	"github.com/nhatminhk63j/uetvoting/server/middleware/authorization"
	"github.com/nhatminhk63j/uetvoting/server/middleware/grpc_error"
	"github.com/nhatminhk63j/uetvoting/services"
)

func Serve(cfg *config.AppConfig) {
	zapSugaredLogger := logger.GetDelegate().(*zap.SugaredLogger)
	zapLogger := zapSugaredLogger.Desugar()
	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)

	grpc_prometheus.EnableHandlingTimeHistogram()

	defer func() {
		sentry.Flush(2 * time.Second)
	}()

	internalServerErr := status.Error(codes.Internal, "Something went wrong in our side.")
	recoveryOpt := grpc_recovery.WithRecoveryHandler(func(err interface{}) error {
		logger.WithFields(logger.Fields{"panic error": err, "stacktrace": string(debug.Stack())}).Error("unexpected error...")
		return internalServerErr
	})

	config := services.NewConfig(cfg.GRPCPort, cfg.HTTPPort)

	e := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDsn,
		Environment: string(cfg.AppMode),
		Debug:       false,
	})
	if e != nil {
		logger.Warnf("sentry.Init: %s", e)
	}

	sv := services.NewServer(config,
		grpcmiddleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger),
			grpcvalidator.UnaryServerInterceptor(),
			grpc_error.UnaryServerInterceptor(cfg.AppMode, internalServerErr),
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
			authentication.UnaryServerInterceptor(jwt.NewJWTResolver()),
			authorization.UnaryServerInterceptor(),
		),
		grpcmiddleware.WithStreamServerChain(
			grpc_prometheus.StreamServerInterceptor,
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(zapLogger),
			grpcvalidator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(recoveryOpt),
		))

	// register handler ...
	if err := sv.Register(
		&health.ServiceServer{},
		auth.InitializeHandler(),
		event.InitializeHandler(),
	); err != nil {
		logger.Fatalf("error register servers: %v", err)
	}

	logger.Infof("serving GRPC: %s", config.GRPC.String())
	logger.Infof("serving HTTP: %s", config.HTTP.String())

	if err := sv.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
