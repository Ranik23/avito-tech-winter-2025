package app

import (
	"avito/config"
	"avito/internal/logger"
	"avito/internal/repository/db/postgres"

	// "avito/internal/repository/db"
	// "avito/internal/repository/db/postgres"
	"avito/internal/router/handlers"
	"avito/internal/service/auth"
	"avito/internal/service/purchase"
	"avito/internal/service/transaction"
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
	authService  		auth.AuthService
	purchaseService 	purchase.PurchaseService
	transactionService 	transaction.TransactionService
	server  	*http.Server
	logger  	*logger.Logger
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

	logger := logger.SetUpLogger(cfg.LogConfig.Level)

	dsn := config.CreatePostgresDSN(cfg)

	storage, err := postgres.NewRepositoryImpl(dsn, logger)
	if err != nil {
		return nil, err
	}

	authService := auth.NewAuthServiceImpl(storage, nil, logger)
	transactionService := transaction.NewTransactionServiceImpl(storage, nil, logger)
	purchaseService := purchase.NewPurchaseServiceImpl(storage, nil, logger)

	router := gin.Default()
	
	routes := router.Group("/api")
	{
		routes.GET("/info", handlers.InfoHandler(nil))
		routes.POST("/sendCoin", handlers.SendCoinHandler(nil))
		routes.GET("/buy/{item}", handlers.BuyHandler(nil))
		routes.POST("/auth", handlers.AuthHandler(authService))
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServerConfig.Host + ":" + cfg.HTTPServerConfig.Port,
		Handler: router,
	}

	return &App{
		server: srv,
		authService: authService,
		transactionService: transactionService,
		purchaseService: purchaseService,
		logger: logger,
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

		select {
		case <-contextShutDown.Done():
			a.logger.Info("Server shutdown complete")
		}
	case err := <-errorCh:
		a.logger.Error("Server encountered an error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *App) ShutDown(ctx context.Context) error {

	a.logger.Info("Shutting Down the server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown the server", slog.String("error", err.Error()))
		return err
	}

	a.logger.Info("Server shutdown successfully")
	return nil
}
