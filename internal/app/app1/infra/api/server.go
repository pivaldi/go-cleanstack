package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/infra/api/gen/cleanstack/v1/cleanstackv1connect"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/infra/api/handler"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/service"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
)

const defaultTimeout = 30 * time.Second

type Server struct {
	port        int
	itemService *service.ItemService
}

func NewServer(port int, itemService *service.ItemService) *Server {
	return &Server{
		port:        port,
		itemService: itemService,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	itemHandler := handler.NewItemHandler(s.itemService)
	path, h := cleanstackv1connect.NewItemServiceHandler(itemHandler)
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
