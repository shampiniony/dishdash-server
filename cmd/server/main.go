package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"dishdash.ru/cmd/server/config"
	httpGateway "dishdash.ru/internal/gateways/http"
	"dishdash.ru/internal/repository/pg"
	"dishdash.ru/internal/usecase"
)

func main() {
	config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	r := httpGateway.NewServer(
		setupUseCases(),
		httpGateway.WithPort(config.C.Port),
	)
	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func setupUseCases() httpGateway.UseCases {
	db, err := pg.NewPostgresDB(context.Background(), pg.Config{
		User:     config.C.PG.User,
		Password: config.C.PG.Password,
		Host:     config.C.PG.Host,
		Port:     config.C.PG.Port,
		Database: config.C.PG.Database,
	})
	if err != nil {
		log.Fatalf("can't setup postgres: %s", err)
	}

	cr := pg.NewCardRepository(db)

	return httpGateway.UseCases{
		Card: usecase.NewCard(cr),
	}
}