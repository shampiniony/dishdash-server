package usecase

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/pkg/filter"
)

type State string

var (
	Swiping  State = "swiping"
	Voting   State = "voting"
	Finished State = "finished"
)

type VoteOption int

var (
	VoteDislike VoteOption
	VoteLike    VoteOption = 1
)

type Match struct {
	ID    int
	Place *domain.Place
}

type Room struct {
	ID   string
	lock sync.RWMutex

	lobby *domain.Lobby

	state      State
	users      map[string]*domain.User
	usersPlace map[string]*domain.Place
	places     []*domain.Place
	swipes     []*domain.Swipe
	matches    []*Match
	matchVotes map[string]VoteOption

	lobbyUseCase Lobby
	placeUseCase Place
}

func NewRoom(
	lobby *domain.Lobby,
	lobbyUseCase Lobby,
	placeUseCase Place,
) *Room {
	return &Room{
		ID:           lobby.ID,
		lobby:        lobby,
		users:        make(map[string]*domain.User),
		places:       make([]*domain.Place, 0),
		usersPlace:   make(map[string]*domain.Place),
		swipes:       make([]*domain.Swipe, 0),
		matchVotes:   make(map[string]VoteOption),
		lobbyUseCase: lobbyUseCase,
		placeUseCase: placeUseCase,
	}
}

func (r *Room) GetNextPlaceForUser(id string) *domain.Place {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.usersPlace[id]
}

func (r *Room) Users() []*domain.User {
	r.lock.RLock()
	defer r.lock.RUnlock()
	users := make([]*domain.User, 0)
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

func (r *Room) Settings() UpdateLobbyInput {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return UpdateLobbyInput{
		ID: r.ID,
		SaveLobbyInput: SaveLobbyInput{
			PriceAvg: r.lobby.PriceAvg,
			Location: r.lobby.Location,
			Tags: filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
				return t.ID
			}),
		},
	}
}

func (r *Room) AddUser(user *domain.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, has := r.users[user.ID]; has {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.users[user.ID] = user
	return nil
}

func (r *Room) RemoveUser(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, has := r.users[id]
	if !has {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.users, id)
	return nil
}

func (r *Room) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.users) == 0
}

func (r *Room) Voting() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == Voting
}

func (r *Room) Swiping() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == Swiping
}

func (r *Room) Finished() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == Finished
}

func (r *Room) UpdateLobby(ctx context.Context, priceAvg int, tagIDs, placeIDs []int64) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.updateLobby(ctx, priceAvg, tagIDs, placeIDs)
}

func (r *Room) updateLobby(ctx context.Context, priceAvg int, tagIDs, placeIDs []int64) error {
	lobby, err := r.lobbyUseCase.UpdateLobby(ctx, UpdateLobbyInput{
		ID: r.lobby.ID,
		SaveLobbyInput: SaveLobbyInput{
			PriceAvg: priceAvg,
			Location: r.lobby.Location,
			Tags:     tagIDs,
			Places:   placeIDs,
		},
	})
	if err != nil {
		return err
	}
	r.lobby = lobby
	return nil
}

func (r *Room) StartSwipes(ctx context.Context) error {
    r.lock.Lock()
    defer r.lock.Unlock()

    log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Request started: Action 'StartSwipes' initiated.", r.lobby.ID, r.ID)

    if len(r.lobby.Tags) == 0 {
        log.Printf("WARNING: [LobbyID: %s] [RoomID: %s] - Warning: Action 'StartSwipes' encountered an issue. Reason: 'No tags found, using default tags'.", r.lobby.ID, r.ID)
        err := r.updateLobby(ctx, 500, []int64{3}, nil)
        if err != nil {
            log.Printf("ERROR: [LobbyID: %s] [RoomID: %s] - Error: Action 'UpdateLobby' failed. Reason: '%v'.", r.lobby.ID, r.ID, err)
            return err
        }
        log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Request successful: Action 'UpdateLobby' completed with default tags.", r.lobby.ID, r.ID)
    }

    var err error
    r.places, err = r.placeUseCase.GetPlacesForLobby(ctx, r.lobby)
    if err != nil {
        log.Printf("ERROR: [LobbyID: %s] [RoomID: %s] - Error: Action 'GetPlacesForLobby' failed. Reason: '%v'.", r.lobby.ID, r.ID, err)
        return err
    }
    log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Request successful: Action 'GetPlacesForLobby' completed successfully.", r.lobby.ID, r.ID)

    err = r.updateLobby(ctx, r.lobby.PriceAvg,
        filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
            return t.ID
        }),
        filter.Map(r.places, func(p *domain.Place) int64 {
            return p.ID
        }))
    if err != nil {
        log.Printf("ERROR: [LobbyID: %s] [RoomID: %s] - Error: Action 'UpdateLobby' failed. Reason: '%v'.", r.lobby.ID, r.ID, err)
        return err
    }
    log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Request successful: Action 'UpdateLobby' completed successfully.", r.lobby.ID, r.ID)

    log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Processing initial place assignments to users.", r.lobby.ID, r.ID)
    for id, _ := range r.users {
        if len(r.places) > 0 {
            r.usersPlace[id] = r.places[0]
            log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - User ID %s assigned to place ID %d.", r.lobby.ID, r.ID, id, r.places[0].ID)
        } else {
            log.Printf("WARNING: [LobbyID: %s] [RoomID: %s] - Warning: No places available to assign to user ID %s.", r.lobby.ID, r.ID, id)
        }
    }

    log.Printf("INFO: [LobbyID: %s] [RoomID: %s] - Swipes successfully started for room.", r.lobby.ID, r.ID)
    return nil
}



func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*Match, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.swipes = append(r.swipes, &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})

	pIdx := slices.IndexFunc(r.places, func(place *domain.Place) bool {
		return place.ID == placeID
	})
	r.usersPlace[userID] = r.places[(pIdx+1)%len(r.places)]

	likes := filter.Count(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})
	if likes > len(r.users)/2 {
		match := &Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
		r.state = Voting
		return match, nil
	}

	return nil, nil
}

func (r *Room) Result() *domain.Place {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.matches[len(r.matches)-1].Place
}

func (r *Room) Vote(userID string, option VoteOption) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.matchVotes[userID] = option
	if len(r.matchVotes) == len(r.users) {
		if r.allVotedLike() {
			r.state = Finished
		} else {
			r.state = Swiping
		}
		r.matchVotes = make(map[string]VoteOption)
	}
}

func (r *Room) allVotedLike() bool {
	for _, vote := range r.matchVotes {
		if vote == VoteDislike {
			return false
		}
	}
	return true
}
