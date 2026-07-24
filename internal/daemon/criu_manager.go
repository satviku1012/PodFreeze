package daemon

import (
	"log/slog"
	"time"
)

// Manager is a minimal placeholder for the CRIU restore workflow.
type Manager struct {
	Logger *slog.Logger
}

func NewManager(logger *slog.Logger) *Manager {
	return &Manager{Logger: logger}
}

// RestoreContainer simulates the restore latency for the initial milestone.
// The real CRIU integration will replace this with an actual restore command.
func (m *Manager) RestoreContainer(containerID string) error {
	if m.Logger != nil {
		m.Logger.Info("restoring container", "container_id", containerID)
	}
	// Buffering logic: pause briefly so the proxy can observe the wake-up window.
	time.Sleep(2 * time.Second)
	return nil
}
