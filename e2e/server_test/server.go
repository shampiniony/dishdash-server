package server_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/cmd/server/config"
	server "dishdash.ru/internal/gateways"
	"dishdash.ru/internal/usecase"
)

func StartServer(cases usecase.Cases) context.CancelFunc {
	config.C.Server.Port = 8081
	ctx, stop := context.WithCancel(context.Background())
	s := server.NewServer(cases)

	go func() {
		log.Info("starting server")
		err := s.Run(ctx)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		log.Info("server closed")
	}()

	for !healthCheck() {
		log.Info("waiting for server to start")
		time.Sleep(1 * time.Second)
	}
	log.Info("server started")

	return stop
}

func healthCheck() bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/api/v1/swagger/index.html", config.C.Server.Port))
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}
