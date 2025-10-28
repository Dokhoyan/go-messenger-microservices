package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Dokhoyan/go-messenger-microservices/api-gateway/internal/config"
)

type GatewayHandler struct {
	userProxyHTTP  *httputil.ReverseProxy
	servicesConfig *config.ServicesConfig
}

func NewGatewayHandler(servicesConfig *config.ServicesConfig) *GatewayHandler {
	return &GatewayHandler{
		servicesConfig: servicesConfig,
		userProxyHTTP:  createProxy(servicesConfig.UserServiceHTTPURL),
	}
}

func createProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Service unavailable: %v", err)))
	}

	// Модифицируем запрос перед проксированием
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Логируем запрос
		log.Printf("Proxying %s %s -> %s", req.Method, req.URL.Path, targetURL)
	}

	return proxy
}

func (h *GatewayHandler) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling auth login: %s %s", r.Method, r.URL.Path)
	// TODO: Add auth service HTTP endpoint via grpc-gateway
	http.Error(w, "Auth service endpoint not implemented yet", http.StatusNotImplemented)
}

func (h *GatewayHandler) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling auth register: %s %s", r.Method, r.URL.Path)
	// TODO: Add auth service HTTP endpoint via grpc-gateway
	http.Error(w, "Auth service endpoint not implemented yet", http.StatusNotImplemented)
}

func (h *GatewayHandler) HandleUserProxy(w http.ResponseWriter, r *http.Request) {
	// Проксируем напрямую к user service HTTP endpoint
	h.userProxyHTTP.ServeHTTP(w, r)
}

func (h *GatewayHandler) HandleChatProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling chat request: %s %s", r.Method, r.URL.Path)
	// TODO: Add chat service HTTP endpoint via grpc-gateway
	http.Error(w, "Chat service endpoint not implemented yet", http.StatusNotImplemented)
}

// HandleStatus handles status requests
func (h *GatewayHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","gateway":"api-gateway","services":{"auth":{"status":"available","url":"%s"},"chat":{"status":"available","url":"%s"},"user":{"status":"available","url":"%s"}}}`,
		h.servicesConfig.AuthServiceURL,
		h.servicesConfig.ChatServiceURL,
		h.servicesConfig.UserServiceURL,
	)
}
