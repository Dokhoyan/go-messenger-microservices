package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dokhoyan/go-messenger-microservices/api-gateway/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/api-gateway/internal/handlers"
)

type App struct {
	config  *config.GatewayConfig
	server  *http.Server
	handler *handlers.GatewayHandler
}

func NewApp(ctx context.Context) (*App, error) {
	gatewayConfig, servicesConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	handler := handlers.NewGatewayHandler(servicesConfig)

	return &App{
		config:  gatewayConfig,
		handler: handler,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/v1/auth/login", a.handler.HandleAuthLogin)
	mux.HandleFunc("/api/v1/auth/register", a.handler.HandleAuthRegister)

	// User routes
	mux.HandleFunc("/api/v1/users/", a.handler.HandleUserProxy)

	// Chat routes
	mux.HandleFunc("/api/v1/chat/", a.handler.HandleChatProxy)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Status endpoint
	mux.HandleFunc("/status", a.handler.HandleStatus)

	a.server = &http.Server{
		Addr:    a.config.Host + ":" + a.config.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("API Gateway starting on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
