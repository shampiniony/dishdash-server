package swipes

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/entities"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(s *socketio.Server, useCases usecase.Cases) {
	s.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	s.OnEvent("/", eventJoinLobby, func(conn socketio.Conn, msg string) {
		var joinEvent joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &joinEvent)
		if err != nil {
			_ = conn.Close()
			return
		}
		_, ok := conn.Context().(*entities.User)
		if ok {
			return
		}

		domLobby, err := useCases.Lobby.GetLobbyByID(context.Background(), joinEvent.LobbyID)
		if err != nil {
			_ = conn.Close()
			return
		}

		lobby, err := entities.FindLobby(domLobby, useCases.Card)
		if err != nil {
			_ = conn.Close()
			return
		}

		u, err := useCases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			_ = conn.Close()
			return
		}

		user := entities.NewUser(*u, useCases.Swipe)
		lobby.Register(conn.ID(), user)

		conn.Join(domLobby.ID)
		conn.SetContext(user)

		s.BroadcastToRoom(
			"",
			user.Lobby.ID,
			eventUserJoined,
			userJoinEvent{
				UserID: user.ID,
				Name:   u.Name,
				Avatar: u.Avatar,
			},
		)

		s.ForEach("/", user.Lobby.ID, func(c socketio.Conn) {
			roomUser, ok := c.Context().(*entities.User)
			if !ok {
				log.Println("Failed to retrieve user from connection context.")
			}
			if roomUser.ID != user.ID {
				u, _ := useCases.User.GetUserByID(context.Background(), roomUser.ID)
				conn.Emit(eventUserJoined, &userJoinEvent{
					UserID: u.ID,
					Name:   u.Name,
					Avatar: u.Avatar,
				})
			}
		})

		settings := user.Lobby.Settings

		conn.Emit(eventSettingsUpdate, &settingsUpdateEvent{
			PriceMin:    settings.PriceMin,
			PriceMax:    settings.PriceMax,
			MaxDistance: settings.MaxDistance,
			Tags:        settings.Tags,
		})
	})

	s.OnEvent("/", eventStartSwipes, func(conn socketio.Conn, msg string) {
		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		s.ForEach("/", user.Lobby.ID, func(c socketio.Conn) {
			roomUser, ok := c.Context().(*entities.User)
			if !ok {
				log.Println("Failed to retrieve user from connection context.")
			}
			if roomUser.ID != user.ID {
				c.Emit(eventStartSwipes)
			}

			firstCard := roomUser.Card()
			c.Emit(eventCard, cardEvent{Card: *firstCard})
		})
	})

	s.OnEvent("/", eventSettingsUpdate, func(conn socketio.Conn, msg string) {
		var updateEvent settingsUpdateEvent
		err := json.Unmarshal([]byte(msg), &updateEvent)
		if err != nil {
			log.Println("wrong settings change event")
			_ = conn.Close()
			return
		}

		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		user.Lobby.UpdateSettings(domain.LobbySettings{
			PriceMin:    updateEvent.PriceMin,
			PriceMax:    updateEvent.PriceMax,
			MaxDistance: updateEvent.MaxDistance,
			Tags:        updateEvent.Tags,
		})

		s.ForEach("/", user.Lobby.ID, func(c socketio.Conn) {
			roomUser, ok := c.Context().(*entities.User)
			if !ok {
				log.Println("Failed to retrieve user from connection context.")
			}

			if roomUser.ID != user.ID {
				c.Emit(eventSettingsUpdate, updateEvent)
			}
		})
	})

	s.OnEvent("/", eventSwipe, func(conn socketio.Conn, msg string) {
		var swipeEvent swipeEvent
		err := json.Unmarshal([]byte(msg), &swipeEvent)
		if err != nil {
			log.Println("wrong swipe event")
			_ = conn.Close()
			return
		}

		u, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		card := u.Card()

		match := u.Swipe(swipeEvent.SwipeType)
		if match != nil {
			s.BroadcastToRoom(
				"",
				u.Lobby.ID,
				eventMatch,
				matchEvent{
					ID:   match.ID,
					Card: *card,
				},
			)

			vote := entities.NewVote(2, func(vote *entities.Vote, results []int) {
				sum := 0
				for _, number := range results {
					sum += number
				}
				if sum == len(u.Lobby.GetUsers()) {
					vote.FinalizeVote()
				}
			}, func(matchResults []int) {
				s.BroadcastToRoom(
					"",
					u.Lobby.ID,
					eventRelaseMatch,
				)

				// 0 - continue
				// 1 - finish
				if matchResults[1] == len(u.Lobby.GetUsers()) {
					lobbyResults := u.Lobby.GetResults()
					if len(lobbyResults) > 1 {
						s.BroadcastToRoom(
							"",
							u.Lobby.ID,
							eventFinalVote,
							&finalVoteEvent{
								Options: lobbyResults,
							},
						)

						vote := entities.NewVote(len(lobbyResults), func(vote *entities.Vote, finalVoteResults []int) {
							sum := 0
							for _, number := range finalVoteResults {
								sum += number
							}
							if sum == len(u.Lobby.GetUsers()) {
								vote.FinalizeVote()
							}
						}, func(results []int) {
							s.BroadcastToRoom(
								"",
								u.Lobby.ID,
								eventFinish,
								&finishEvent{
									Result: lobbyResults[0],
								},
							)
						})

						time.AfterFunc(20*time.Second, func() {
							vote.FinalizeVote()
						})

					} else {
						s.BroadcastToRoom(
							"",
							u.Lobby.ID,
							eventFinish,
							&finishEvent{
								Result: lobbyResults[0],
							},
						)
					}
				}
			})

			u.Lobby.RegisterVote(vote, match.ID)
		}

		newCard := u.Card()

		conn.Emit(eventCard, cardEvent{Card: *newCard})
	})

	s.OnEvent("/", eventVote, func(conn socketio.Conn, msg string) {
		var voteEvent voteEvent
		err := json.Unmarshal([]byte(msg), &voteEvent)
		if err != nil {
			log.Println("wrong swipe event")
			_ = conn.Close()
			return
		}

		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		v := user.Lobby.GetVoteByID(voteEvent.VoteID)
		v.Vote(int(voteEvent.VoteOption))

		u, _ := useCases.User.GetUserByID(context.Background(), user.ID)
		s.BroadcastToRoom(
			"",
			user.Lobby.ID,
			eventVoted,
			&votedEvent{
				User: User{
					UserID: u.ID,
					Name:   u.Name,
					Avatar: u.Avatar,
				},
				VoteID:     voteEvent.VoteID,
				VoteOption: voteEvent.VoteOption,
			},
		)
	})

	s.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		user.Lobby.Unregister(conn.ID())

		u, _ := useCases.User.GetUserByID(context.Background(), user.ID)

		s.BroadcastToRoom("", user.Lobby.ID, eventUserLeft, &userJoinEvent{
			UserID: u.ID,
			Name:   u.Name,
			Avatar: u.Avatar,
		})
	})
}
