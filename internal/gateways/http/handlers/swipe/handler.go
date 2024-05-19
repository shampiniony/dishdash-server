package swipe

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"dishdash.ru/internal/usecase/swipe"

	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

func SetupLobby(wsServer *socketio.Server, useCases usecase.Cases) {
	wsServer.OnConnect("", func(conn socketio.Conn) error {
		conn.SetContext("")
		log.Println("connected:", conn.ID())
		return nil
	})

	wsServer.OnEvent("", "echo", func(s socketio.Conn, msg string) {
		log.Println("echo:", msg)
		s.Emit("echo", msg)
	})

	wsServer.OnEvent("", eventJoinLobby, func(conn socketio.Conn, msg string) {
		var e joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &e)
		if err != nil {
			_ = conn.Close()
			return
		}
		_, ok := conn.Context().(*swipe.User)
		if ok {
			return
		}

		domLobby, err := useCases.Lobby.GetLobbyByID(context.Background(), e.LobbyID)
		if err != nil {
			_ = conn.Close()
			return
		}

		lobby, err := swipe.FindLobby(domLobby, useCases.Card)
		if err != nil {
			_ = conn.Close()
			return
		}

		u := swipe.NewUser(conn.ID(), useCases.Swipe)

		conn.Join(strconv.FormatInt(domLobby.ID, 10))
		lobby.Register(u)
		conn.SetContext(u)

		firstCard := u.Card()
		conn.Emit(eventCard, cardEvent{Card: firstCard.ToDto()})
	})

	wsServer.OnEvent("", eventSwipe, func(conn socketio.Conn, msg string) {
		var swipeEvent swipeEvent
		err := json.Unmarshal([]byte(msg), &swipeEvent)
		if err != nil {
			log.Println("wrong swipe event")
			_ = conn.Close()
			return
		}

		u, ok := conn.Context().(*swipe.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		card := u.Card()
		match := u.Swipe(swipeEvent.SwipeType)
		if match != nil {
			wsServer.BroadcastToRoom(
				"",
				strconv.FormatInt(u.Lobby.Id, 10),
				eventMatch,
				matchEvent{
					Card: card.ToDto(),
				},
			)
		}

		newCard := u.Card()
		conn.Emit(eventCard, cardEvent{Card: newCard.ToDto()})
	})

	wsServer.OnError("", func(c socketio.Conn, err error) {
		log.Println("socket.io error: ", err)
	})

	wsServer.OnDisconnect("", func(s socketio.Conn, reason string) {
		u, ok := s.Context().(*swipe.User)
		if ok {
			u.Lobby.Unregister(u)
		}
		log.Println("disconnected:", s.ID(), "reason:", reason)
	})
}
