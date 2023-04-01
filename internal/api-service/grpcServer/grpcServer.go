package grpcServer

import (
	"context"
	"fmt"
	"log"

	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"google.golang.org/grpc"
)

type grpcServer struct {
	loginService.UnimplementedLoginServiceServer

	config GrpcConfig
	logger logger.ILogger
	conn   *grpc.ClientConn
	client loginService.LoginServiceClient
}

func NewGrpcServer(config GrpcConfig, logger logger.ILogger) *grpcServer {
	return &grpcServer{
		logger: logger,
		config: config,
	}
}

func (g *grpcServer) Run() {
	var opts []grpc.DialOption
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), g.config.ConnectionTimeout)
	defer cancel()

	g.conn, err = grpc.DialContext(ctx, g.config.Addr, opts...)
	if err != nil {
		log.Panic("grpc.DialContext", err)
	}

	g.client = loginService.NewLoginServiceClient(g.conn)
}

func (g *grpcServer) VerifyUser(ctx context.Context, loginreq *loginService.LoginRequest) (res *loginService.LoginResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.config.FunctionCallTimeout)
	defer cancel()

	g.logger.Debug("grpc.VerifyUser", fmt.Sprintf("{username: %s, password: %s}\n", loginreq.Username, loginreq.Password))
	res, err = g.client.VerifyUser(ctx, loginreq)
	g.logger.Debug("grpc.VerifyUser", fmt.Sprintf("response: %s", res))
	return res, err
}

func (g *grpcServer) Close() {
	g.conn.Close()
}
