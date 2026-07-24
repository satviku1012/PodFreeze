package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Server buffers inbound HTTP traffic until the orchestrator confirms the container is ready.
type Server struct {
	Addr   string
	Logger *slog.Logger
	Client *Client
}

func NewServer(addr string, logger *slog.Logger, client *Client) *Server {
	return &Server{Addr: addr, Logger: logger, Client: client}
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	s.Logger.Info("proxy listening", "addr", s.Addr)
	return http.ListenAndServe(s.Addr, mux)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	containerID := r.Header.Get("X-Container-ID")
	if containerID == "" {
		containerID = strings.TrimPrefix(r.URL.Path, "/")
	}
	targetIP := r.Header.Get("X-Target-IP")
	if targetIP == "" {
		targetIP = "127.0.0.1:8080"
	}

	if containerID == "" {
		http.Error(w, "missing container id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	if s.Logger != nil {
		s.Logger.Info("buffering request", "container_id", containerID, "target_ip", targetIP)
	}

	resp, err := s.Client.WakeContainer(ctx, containerID, targetIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("wake failed: %v", err), http.StatusBadGateway)
		return
	}
	if !resp.Success {
		http.Error(w, resp.Message, http.StatusBadGateway)
		return
	}

	proxyTarget, err := url.Parse("http://" + targetIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid target: %v", err), http.StatusBadRequest)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyTarget)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, fmt.Sprintf("proxy error: %v", err), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}
