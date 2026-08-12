package adapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arinamklvch/xkcd-helper/internal/db"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersStorage struct {
	queries *db.Queries
}

func NewUsersStorage(pool *pgxpool.Pool) *UsersStorage {
	return &UsersStorage{
		queries: db.New(pool),
	}
}

var ErrUserNotFound = errors.New("user not found")

func (u *UsersStorage) GetUser(login, password string) (domain.User, error) {
	dbUser, err := u.queries.GetUser(context.Background(), db.GetUserParams{
		Login:    pgtype.Text{String: login, Valid: true},
		Password: pgtype.Text{String: password, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		Login:    dbUser.Login.String,
		Password: dbUser.Password.String,
		Role:     dbUser.Role.String,
	}, nil
}
