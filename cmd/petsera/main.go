package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vliubezny/petsera/internal/config"
	"github.com/vliubezny/petsera/internal/health"
	"github.com/vliubezny/petsera/internal/server"
	"github.com/vliubezny/petsera/internal/storage/gcs"
	"github.com/vliubezny/petsera/internal/storage/postgres"
	"github.com/vliubezny/petsera/ui"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg := config.MustLoad()

	ctx := context.Background()

	pgStorage := postgres.MustSetupStorage(postgres.Config{
		Host:            cfg.DBHost,
		Port:            cfg.DBPort,
		User:            cfg.DBUser,
		Password:        cfg.DBPassword,
		DBName:          cfg.DBName,
		MaxOpenConns:    cfg.DBMaxOpenConnections,
		MaxIdleConns:    cfg.DBMaxIdleConnections,
		EnableMigration: cfg.DBEnableMigration,
		MigrationsDir:   cfg.DBMigrations,
	})
	defer pgStorage.Close()

	fileStorage, err := gcs.NewGCS(ctx, cfg.GCSBucket)
	if err != nil {
		logrus.WithError(err).Fatal("failed to init file storage")
	}

	checker := health.SetupChecks(pgStorage, fileStorage)

	assets, err := ui.LoadFileSystem()
	if err != nil {
		logrus.WithError(err).Fatal("failed to access embedded file system")
	}

	srv, err := server.New(server.Config{
		DevMode:             cfg.DevMode,
		Port:                cfg.HTTPPort,
		Statics:             assets,
		AnnouncementStorage: pgStorage,
		FileStorage:         fileStorage,
		Checker:             checker,
		FrontendConfig: map[string]any{
			"apiKey": cfg.GoogleAPIKey,
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to create server")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(ctx); err != nil {
			logrus.WithError(err).Fatal("failed to run server")
		}
	}()

	<-sigs
	logrus.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to stop gracefully")
	}
}
