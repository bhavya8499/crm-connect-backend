package grpc

import (
	"connect-crm-backend/crm-connect-backend/internal/config"
	"connect-crm-backend/crm-connect-backend/internal/constant"
	interceptors "connect-crm-backend/crm-connect-backend/internal/protocol/grpc/interceptors/auth"
	cgrpc "github.com/swiggy-private/gocommons/grpc"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
)

var (
	grpcServer *Server
	onceInit   *sync.Once
)

type Server struct {
	server           *grpc.Server
	middlewareConfig *config.MiddlewareConfig
	config           *config.Config
}

func NewServer(server *grpc.Server) *Server {

	onceInit.Do(func() {
		authInterceptor := interceptors.NewAuthInterceptor(config.NewMiddlewareConfig())
		serverOption := grpc.ChainUnaryInterceptor(authInterceptor.Unary())
		gg := cgrpc.WithGrpcServerOption(serverOption)
		grpcServer = &Server{
			server: cgrpc.NewServer(gg),
		}
	})

	//register client/stub here
	return grpcServer
}

func (s *Server) Start() {

	listener, err := net.Listen(constant.TCP, ":"+s.config.ServerConfig.GRPCPort)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.server.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) ShutDownGracefully(sgnl os.Signal) {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
