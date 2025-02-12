package app

import (
	"avito/config"
	"avito/internal/logger"
	"avito/internal/router/handlers"
	"avito/internal/usecase"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)


type App struct {
	UseCase 	usecase.UserCase
	Server 		*http.Server
	Logger 		*logger.Logger
	Context 	context.Context
}

func NewApp(configPath string) (*App, error) {

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg config.Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	router := gin.Default()

	routes := router.Group("/api")
	{
		routes.GET("/info", handlers.InfoHandler(nil))
		routes.POST("/sendCoin", handlers.SendCoinHandler(nil))
		routes.GET("/buy/{item}", handlers.BuyHandler(nil))
		routes.POST("/auth", handlers.AuthHandler(nil))
	}

	srv := &http.Server{
		Addr: cfg.HTTPServerConfig.Host + cfg.HTTPServerConfig.Port,
		Handler: router,
	}

	return &App{
		Server: srv,
	}, nil
}

func (a *App) Run() error {

	errorCh := make(chan error)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("Failed to start the server", slog.String("error", err.Error()))
			errorCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		contextShutDown, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownRelease()

		if err := a.ShutDown(contextShutDown); err != nil {
			a.Logger.Error("Failed to shutdown the server", slog.String("err", err.Error()))
			return err
		}

		select {
		case <-contextShutDown.Done():
			a.Logger.Info("Server shutdown complete")
		}
	case err := <-errorCh:
		a.Logger.Error("Server encountered an error", slog.String("error", err.Error()))
		return err
	}

	return nil
}



func (a *App) ShutDown(ctx context.Context) error {

	a.Logger.Info("Shutting Down the server")

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("Failed to shutdown the server", slog.String("error", err.Error()))
		return err
	}

	a.Logger.Info("Server shutdown successfully")
	return nil
}