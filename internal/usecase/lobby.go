package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type LobbyUseCase struct {
	lobbyRepo repo.Lobby
}

func NewLobbyUseCase(lobbyRepo repo.Lobby) *LobbyUseCase {
	return &LobbyUseCase{lobbyRepo: lobbyRepo}
}

func (l *LobbyUseCase) CreateLobby(ctx context.Context, lobbyInput LobbyInput) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		Location: lobbyInput.Location,
	}
	lobby, err := l.lobbyRepo.CreateLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}
	return lobby, err
}

func (l *LobbyUseCase) NearestActiveLobby(ctx context.Context, loc domain.Coordinate) (*domain.Lobby, float64, error) {
	return l.lobbyRepo.NearestActiveLobby(ctx, loc)
}

func (l *LobbyUseCase) DeleteLobbyByID(ctx context.Context, id string) error {
	return l.lobbyRepo.DeleteLobbyByID(ctx, id)
}

func (l *LobbyUseCase) GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error) {
	return l.lobbyRepo.GetLobbyByID(ctx, id)
}

func (l *LobbyUseCase) GetCardsForSettings(ctx context.Context, loc domain.Coordinate, settings *domain.LobbySettings) ([]*domain.Card, error) {
	return l.lobbyRepo.GetCardsForSettings(ctx, loc, settings)
}

func (l *LobbyUseCase) SetLobbyActive(ctx context.Context, id string, active bool) error {
	return l.lobbyRepo.SetLobbyActive(ctx, id, active)
}
