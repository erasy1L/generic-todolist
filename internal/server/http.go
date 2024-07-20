package server

import (
	"context"
	"io/fs"
	"net/http"
	"os"

	"todo-list/internal/service/list"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

var FS fs.FS

type HTTPServer struct {
	Port string

	taskService list.Service

	server *http.Server
}

func NewHTTPServer(port string, taskService list.Service) *HTTPServer {
	return &HTTPServer{
		Port:        port,
		taskService: taskService,
	}
}

func (s *HTTPServer) Run() error {
	router := chi.NewRouter()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &logger}))
	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

	router.Handle("/api/docs/*", http.StripPrefix("/api/docs", http.FileServer(http.FS(FS))))

	router.Mount("/api", HandlerFromMux(s, router))

	s.server = &http.Server{
		Addr:    ":" + s.Port,
		Handler: router,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("failed to start server")
		}
	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}

	return nil
}
