package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	pb "github.com/podfreeze/core/api/proto/podfreeze"
)

// Server exposes the orchestrator RPC interface for the proxy.
type Server struct {
	pb.UnimplementedOrchestratorServer
	Manager *Manager
	Logger  *slog.Logger
}

func (s *Server) WakeContainer(ctx context.Context, req *pb.WakeRequest) (*pb.WakeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("wake request is nil")
	}
	if s.Logger != nil {
		s.Logger.Info("wake request received", "container_id", req.ContainerId, "target_ip", req.TargetIp)
	}

	if err := s.Manager.RestoreContainer(req.ContainerId); err != nil {
		return &pb.WakeResponse{Success: false, Message: err.Error(), ReadyTimeMs: 0}, nil
	}

	return &pb.WakeResponse{Success: true, Message: "container restored", ReadyTimeMs: 2000}, nil
}

func ListenAndServe(addr string, logger *slog.Logger) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	manager := NewManager(logger)
	pb.RegisterOrchestratorServer(grpcServer, &Server{Manager: manager, Logger: logger})

	logger.Info("grpc daemon listening", "addr", addr)
	return grpcServer.Serve(listener)
}
