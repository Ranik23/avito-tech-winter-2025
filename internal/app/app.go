package app

import (
	"avito/config"
	"avito/internal/logger"
	"avito/internal/repository/postgres"
	"avito/internal/service"
	"avito/internal/router/handlers"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
)

type App struct {
	serviceManager 	*service.ServiceManager
	server  		*http.Server
	logger  		*logger.Logger
	cfg				*config.Config
}

func NewApp(configPath string) (*App, error) {

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger(cfg.LogConfig.Level)

	dsn := config.CreatePostgresDSN(cfg)

	storage, err := postgres.NewRepositoryImpl(dsn, logger)
	if err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		return nil, err
	}

	serviceManager := service.NewServiceManager(storage, logger)

	router := gin.Default()
	
	routes := router.Group("/api")
	{
		routes.GET("/info", handlers.InfoHandler(serviceManager))
		routes.POST("/sendCoin", handlers.SendCoinHandler(serviceManager))
		routes.GET("/buy/:item", handlers.BuyHandler(serviceManager))
		routes.POST("/auth", handlers.AuthHandler(serviceManager))
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServerConfig.Host + ":" + cfg.HTTPServerConfig.Port,
		Handler: router,
	}

	logger.Info("Application initialized successfully")

	return &App{
		server: srv,
		cfg: cfg,
		logger: logger,
		serviceManager: serviceManager,
	}, nil
}

func (a *App) Run() error {

	errorCh := make(chan error)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Failed to start the server", slog.String("error", err.Error()))
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
			a.logger.Error("Failed to shutdown the server", slog.String("err", err.Error()))
			return err
		}

		<-contextShutDown.Done()
		a.logger.Info("Server shutdown complete")
	case err := <-errorCh:
		a.logger.Error("Server encountered an error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *App) ShutDown(ctx context.Context) error {

	a.logger.Info("Shutting Down the server...")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown the server", slog.String("error", err.Error()))
		return err
	}

	a.logger.Info("Server shutdown successfully")
	return nil
}
