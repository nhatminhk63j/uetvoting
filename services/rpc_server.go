package services

import (
	"context"
	"fmt"
	"github.com/nhatminhk63j/uetvoting/pb/auth/v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/nhatminhk63j/uetvoting/pb/health/v1"
)

// DefaultConfig return a default server config
func DefaultConfig() Config {
	return NewConfig(10443, 10080)
}

// NewConfig return a optional config with grpc port and http port.
func NewConfig(grpcPort, httpPort int) Config {
	return Config{
		GRPC: ServerListen{
			Host: "0.0.0.0",
			Port: grpcPort,
		},
		HTTP: ServerListen{
			Host: "0.0.0.0",
			Port: httpPort,
		},
	}
}

// Config hold http/grpc server config
type Config struct {
	GRPC ServerListen `json:"grpc" mapstructure:"grpc" yaml:"grpc"`
	HTTP ServerListen `json:"http" mapstructure:"http" yaml:"http"`
}

// ServerListen config for host/port socket listener
type ServerListen struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

// String return socket listen DSN
func (l ServerListen) String() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

// Server structure
type Server struct {
	gRPC *grpc.Server
	mux  *runtime.ServeMux
	cfg  Config
}

func NewServer(cfg Config, opt ...grpc.ServerOption) *Server {
	return &Server{
		gRPC: grpc.NewServer(opt...),
		mux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}),
		),
		cfg: cfg,
	}
}

// Register ...
func (s *Server) Register(grpcServer ...interface{}) error {
	for _, srv := range grpcServer {
		switch _srv := srv.(type) {
		case health.HealthCheckServiceServer:
			health.RegisterHealthCheckServiceServer(s.gRPC, _srv)
			if err := health.RegisterHealthCheckServiceHandlerFromEndpoint(context.Background(), s.mux, s.cfg.GRPC.String(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
				return err
			}
		case auth.AuthServiceServer:
			auth.RegisterAuthServiceServer(s.gRPC, _srv)
			if err := auth.RegisterAuthServiceHandlerFromEndpoint(context.Background(), s.mux, s.cfg.GRPC.String(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown GRPC Service to register %#v", srv)
		}
	}
	return nil
}

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	glog.Infof("preflight request for %s", r.URL.Path)
}

// Serve server listen for HTTP and GRPC
func (s *Server) Serve() error {
	stop := make(chan os.Signal, 1)
	errch := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", s.mux)
	httpServer := http.Server{
		Addr:    s.cfg.HTTP.String(),
		Handler: allowCORS(httpMux),
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errch <- err
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", s.cfg.GRPC.String())
		if err != nil {
			errch <- err
			return
		}
		if err := s.gRPC.Serve(listener); err != nil {
			errch <- err
		}
	}()
	for {
		select {
		case <-stop:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			s.gRPC.GracefulStop()
			httpServer.Shutdown(ctx)
			if !isRunningInDockerContainer() {
				fmt.Println("Shutting down. Wait for 1 second")
				time.Sleep(1 * time.Second)
			}
			return nil
		case err := <-errch:
			return err
		}
	}
}
