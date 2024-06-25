package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	ws "dishdash.ru/internal/gateways/ws"
	http "dishdash.ru/internal/gateways/http"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"

	"github.com/tj/go-spin"
	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer *http.Server
	WsServer   *ws.Server
}

func NewServer(useCases usecase.Cases) *Server {

	s := &Server{
		HttpServer: http.NewServer(useCases),
		WsServer: ws.NewServer(useCases),
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.HttpServer.ListenAndServe()
	})
	eg.Go(func() error {
		return s.wsServer.Serve()
	})

	<-ctx.Done()
	err := s.HttpServer.Shutdown(ctx)
	err = errors.Join(err, s.WsServer.Close())
	err = errors.Join(eg.Wait(), err)
	shutdownWait()
	return err
}

func newSocketIOServer() *socketio.Server {
	wt := websocket.Default
	// TODO legal CheckOrigin
	wt.CheckOrigin = func(_ *http.Request) bool {
		return true
	}

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			wt,
		},
	})

	return server
}

func shutdownWait() {
	spinner := spin.New()
	const spinIterations = 20
	for range spinIterations {
		fmt.Printf("\rgraceful shutdown %s ", spinner.Next())
		time.Sleep(shutdownDuration / spinIterations)
	}
	fmt.Println()
}
