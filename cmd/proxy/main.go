package main

import (
	"os"

	"github.com/podfreeze/core/internal/proxy"
	"github.com/podfreeze/core/pkg/logger"
)

func main() {
	log := logger.New()
	client, err := proxy.NewClient("127.0.0.1:50051", log)
	if err != nil {
		log.Error("failed to initialize grpc client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	server := proxy.NewServer(":8081", log, client)
	log.Info("starting podfreeze proxy")
	if err := server.ListenAndServe(); err != nil {
		log.Error("proxy failed", "error", err)
		os.Exit(1)
	}
}
