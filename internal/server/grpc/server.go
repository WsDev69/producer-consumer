package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/WsDev69/producer-consumer/internal/config"
	pb "github.com/WsDev69/producer-consumer/internal/handler/grpc/gen/task"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	grpc *grpc.Server
	cfg  config.GRPC
	h    pb.TaskServerServer
}

func NewServer(cfg config.GRPC, h pb.TaskServerServer) *Server {
	return &Server{cfg: cfg, h: h}
}

func (s *Server) Serve(ctx context.Context, wg *sync.WaitGroup, opt ...grpc.ServerOption) {
	wg.Add(1)
	go func() {
		log := slog.Default()
		if s.cfg.TLS {
			creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFile, s.cfg.KeyFile)
			if err != nil {
				log.Error(fmt.Sprintf("failed to generate credentials: %v", err))
				os.Exit(1)
			}
			opt = append(opt, grpc.Creds(creds))
		}

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
		if err != nil {
			log.Error(fmt.Sprintf("failed to listen: %v", err))
			os.Exit(1)
		}

		s.grpc = grpc.NewServer(opt...)
		pb.RegisterTaskServerServer(s.grpc, s.h)
		slog.Default().Info(fmt.Sprintf("grpc server listening on %s:%d", s.cfg.Host, s.cfg.Port))
		if err := s.grpc.Serve(lis); err != nil {
			log.Error(fmt.Sprintf("failed to serve: %v", err))
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Default().Warn("shutting down grpc server")
		if s.grpc != nil {
			s.grpc.GracefulStop()
		}

		slog.Default().Warn("grpc server stopped")
	}()
}
