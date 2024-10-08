package tests

import (
	"testing"
	"time"

	"dishdash.ru/e2e/pg_test"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func LobbyVote(t *testing.T) *SocketIOSession {
	user1 := postUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	user2 := postUserWithID(t, &domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})

	lobby := findLobby(t)

	cli1, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	cli2, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	sioSess := newSocketIOSession()
	sioSess.addUser(user1.Name)
	sioSess.addUser(user2.Name)

	listenEvent := func(eventName string) {
		cli1.OnEvent(eventName, sioSess.sioAddFunc(user1.Name, eventName))
		cli2.OnEvent(eventName, sioSess.sioAddFunc(user2.Name, eventName))
	}

	listenEvent(event.Match)
	listenEvent(event.Voted)
	listenEvent(event.ReleaseMatch)
	listenEvent(event.Finish)

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	cli1Emit := emitWithLogFunc(cli1, user1.Name)
	cli2Emit := emitWithLogFunc(cli2, user2.Name)

	sioSess.newStep("Joining lobby")
	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	cli1Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{pg_test.Tags[3].ID},
	})
	time.Sleep(waitTime)

	sioSess.newStep("Start swipes")
	cli1Emit(event.StartSwipes)
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes (1)")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Vote like and dislike")
	cli1Emit(event.Vote, event.VoteEvent{ID: 0, Option: usecase.VoteLike})
	cli2Emit(event.Vote, event.VoteEvent{ID: 0, Option: usecase.VoteDislike})
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes (2)")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Vote both likes")
	cli1Emit(event.Vote, event.VoteEvent{ID: 1, Option: usecase.VoteLike})
	cli2Emit(event.Vote, event.VoteEvent{ID: 1, Option: usecase.VoteLike})
	time.Sleep(waitTime)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}
