package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pivaldi/go-cleanstack/internal/app/user/api/gen/user/v1/userv1connect"
	"github.com/pivaldi/go-cleanstack/internal/app/user/api/handler"
	"github.com/pivaldi/go-cleanstack/internal/app/user/service"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
)

const defaultTimeout = 30 * time.Second

type Server struct {
	port        int
	userService *service.UserService
}

func NewServer(port int, userService *service.UserService) *Server {
	return &Server{
		port:        port,
		userService: userService,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	userHandler := handler.NewUserHandler(s.userService)
	path, h := userv1connect.NewUserServiceHandler(userHandler)
	mux.Handle(path, h)

	addr := fmt.Sprintf(":%d", s.port)
	logging.GetLogger().Info("starting HTTP server", logging.String("address", addr))

	h2server := &http2.Server{
		IdleTimeout:      defaultTimeout,
		WriteByteTimeout: defaultTimeout,
	}

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(mux, h2server),
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		IdleTimeout:  defaultTimeout,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
