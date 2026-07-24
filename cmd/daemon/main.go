package main

import (
	"os"

	"github.com/podfreeze/core/internal/daemon"
	"github.com/podfreeze/core/pkg/logger"
)

func main() {
	log := logger.New()
	log.Info("starting podfreeze daemon")
	if err := daemon.ListenAndServe(":50051", log); err != nil {
		log.Error("daemon failed", "error", err)
		os.Exit(1)
	}
}
