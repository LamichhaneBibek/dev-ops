package apiserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/LamichhaneBibek/dev-ops/config"
	"github.com/LamichhaneBibek/dev-ops/store"
)

type ApiServer struct {
	config     *config.Config
	logger     *slog.Logger
	store      *store.Store
	jwtManager *JWTManager
}

func New(config *config.Config, logger *slog.Logger, store *store.Store, jwtManager *JWTManager) *ApiServer {
	return &ApiServer{
		config:     config,
		logger:     logger,
		store:      store,
		jwtManager: jwtManager,
	}
}

func (s *ApiServer) welcome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Golang api server"))
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.welcome)
	mux.HandleFunc("GET /ping", s.ping)
	mux.HandleFunc("POST /auth/signup", s.signupHandler())
	mux.HandleFunc("POST /auth/signin", s.signinHandler())

	middleware := NewLoggerMiddleware(s.logger)

	srv := &http.Server{
		Addr:    net.JoinHostPort(s.config.ApiserverHost, s.config.ApiserverPort),
		Handler: middleware(mux),
	}

	go func() {
		s.logger.Info("apiserver running", "port", s.config.ApiserverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("apiserver failed to listen and server", "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("apiserver failed to shutdown", "error", err)
		}
	}()

	wg.Wait()
	return srv.ListenAndServe()
}
