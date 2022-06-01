package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vliubezny/petsera/internal/config"
	"github.com/vliubezny/petsera/internal/server"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()

	// docsDB, err := postgres.New(postgres.Config{
	// 	Host:     cfg.DBHost,
	// 	Port:     cfg.DBPort,
	// 	User:     cfg.DBUser,
	// 	Password: cfg.DBPassword,
	// 	DBName:   cfg.DBName,
	// })
	// if err != nil {
	// 	log.Fatalf("failed to init doc DB: %v", err)
	// }

	// defer docsDB.Close()

	// assets, err := ui.LoadFileSystem()
	// if err != nil {
	// 	log.Fatalf("failed to access embedded file system: %v", err)
	// }

	srv, err := server.New(server.Config{
		Port:        cfg.HTTPPort,
		// Statics:     assets,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	<-sigs
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop gracefully: %v", err)
	}
}
