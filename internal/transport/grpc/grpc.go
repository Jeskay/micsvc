package transport

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Jeskay/micsvc/internal/auth"
	"github.com/Jeskay/micsvc/internal/transport/grpc/interceptors"
	"github.com/Jeskay/micsvc/internal/transport/grpc/server"
	"github.com/Jeskay/micsvc/internal/user"
	proto "github.com/Jeskay/micsvc/protos"
)

type RPCServer struct {
	userSvc    *user.Service
	authSvc    *auth.Service
	baseServer *grpc.Server
}

func NewRPCServer(userSvc *user.Service, authSvc *auth.Service) *RPCServer {
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.NewAuthUnary(authSvc)))
	userServer := server.NewUserGRPC(authSvc, userSvc)

	reflection.Register(baseServer)
	proto.RegisterUsersServer(baseServer, userServer)

	return &RPCServer{
		userSvc:    userSvc,
		authSvc:    authSvc,
		baseServer: baseServer,
	}
}

func (s *RPCServer) Run(addr string) error {
	grpcListener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.baseServer.Serve(grpcListener)
}

func (s *RPCServer) Shutdown(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		s.baseServer.GracefulStop()
		close(done)
	}()	
	select {
	case <-done:
		log.Println("RPC Server stopped gracefully")
	case <-ctx.Done():
		log.Println("Failed to stop gracefully")
		s.baseServer.Stop()
	}
}
