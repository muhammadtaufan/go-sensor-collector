package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (types.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) (UserRepository, error) {
	return &userRepository{
		db: db,
	}, nil
}

func (ur *userRepository) GetUser(ctx context.Context, username string) (types.User, error) {
	query := `SELECT id, username, password from users WHERE username = ?`

	var userData types.User
	err := ur.db.Get(&userData, query, username)
	if err != nil {
		fmt.Println(err)
		return types.User{}, err
	}

	fmt.Println("masujk")
	return userData, nil
}
