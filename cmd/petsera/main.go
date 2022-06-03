package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vliubezny/petsera/internal/config"
	"github.com/vliubezny/petsera/internal/server"
	"github.com/vliubezny/petsera/internal/storage/postgres"
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

	// assets, err := ui.LoadFileSystem()
	// if err != nil {
	// 	log.Fatalf("failed to access embedded file system: %v", err)
	// }

	srv, err := server.New(server.Config{
		Port: cfg.HTTPPort,
		// Statics:     assets,
		AnnouncementStorage: pgStorage,
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
