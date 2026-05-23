package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doter2410/avito-project/internal/core/transport/transport"
	"github.com/doter2410/avito-project/internal/courier"
	"github.com/doter2410/avito-project/internal/storage/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "7540"
	}
	port := pflag.String("port", envPort, "Server port")
	pflag.Parse()

	logI := log.New(os.Stdout, "Server", log.LstdFlags)

	psUser := os.Getenv("POSTGRES_USER")
	psPass := os.Getenv("POSTGRES_PASSWORD")
	psDb := os.Getenv("POSTGRES_DB")
	psPort := os.Getenv("POSTGRES_PORT")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", psUser, psPass, psPort, psDb)
	psCon, err := postgres.New(ctx, connectionString)
	if err != nil {
		logI.Fatal(err)
	}
	defer psCon.Close()

	storage := courier.NewStorage(psCon)

	handler := courier.NewHandler(storage)

	srv := transport.NewServer(port, logI, handler)

	go func() {
		err := srv.HttpServer.ListenAndServe()

		if err != nil {
			logI.Fatal(err)
		}

	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.HttpServer.Shutdown(shutdownCtx)

}
