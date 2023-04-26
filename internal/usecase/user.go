package usecase

import (
	"context"

	"github.com/muhammadtaufan/go-sensor-collector/internal/repository"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type User interface {
	GetUser(ctx context.Context, username string) (types.User, error)
}

type user struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) User {
	return &user{repo: repo}
}

func (uu *user) GetUser(ctx context.Context, username string) (types.User, error) {
	results, err := uu.repo.GetUser(ctx, username)
	if err != nil {
		return types.User{}, err
	}
	return results, nil
}
