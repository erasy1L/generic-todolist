package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-list/config"
	"todo-list/internal/repository"
	"todo-list/internal/repository/postgres"
	"todo-list/internal/server"
	"todo-list/internal/service/list"
	"todo-list/log"

	"github.com/rs/zerolog"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
		return
	}

	db, err := postgres.New(configs.DB)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		return
	}
	defer db.Close()

	err = db.Migrate()
	if err != nil {
		panic(err)
	}

	httpServer, err := setupServer(&configs, &db, logger)
	if err != nil {
		return
	}

	runServer(httpServer, &configs, logger)

}

func setupServer(configs *config.Configs, db *postgres.DB, logger *zerolog.Logger) (*server.HTTPServer, error) {
	taskRepository := repository.MustNew(logger, db.Client)

	taskService := list.New(list.WithTaskRepository(taskRepository))

	httpServer := server.NewHTTPServer(configs.APP.Port, taskService)
	if err := httpServer.Run(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
		return nil, err
	}

	return httpServer, nil
}

func runServer(httpServer *server.HTTPServer, configs *config.Configs, logger *zerolog.Logger) {
	fmt.Println("server started at http://localhost:" + configs.APP.Port +
		", swagger available at http://localhost:8080/api/docs/index.htm")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("failed to stop server")
		return
	}

	fmt.Println("running cleanup tasks...")

	fmt.Println("server stopped successfully")
}
