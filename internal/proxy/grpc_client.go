package proxy

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/podfreeze/core/api/proto/podfreeze"
)

// Client talks to the local daemon over gRPC.
type Client struct {
	conn   *grpc.ClientConn
	client pb.OrchestratorClient
	logger *slog.Logger
}

func NewClient(addr string, logger *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to daemon: %w", err)
	}
	return &Client{conn: conn, client: pb.NewOrchestratorClient(conn), logger: logger}, nil
}

func (c *Client) WakeContainer(ctx context.Context, containerID, targetIP string) (*pb.WakeResponse, error) {
	if c.client == nil {
		return nil, fmt.Errorf("grpc client not initialized")
	}
	resp, err := c.client.WakeContainer(ctx, &pb.WakeRequest{ContainerId: containerID, TargetIp: targetIP})
	if err != nil {
		return nil, fmt.Errorf("wake container rpc: %w", err)
	}
	return resp, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
