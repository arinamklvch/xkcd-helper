package adapter

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/arinamklvch/xkcd-helper/internal/db"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersStorage struct {
	queries *db.Queries
	logger  *slog.Logger
}

func NewUsersStorage(pool *pgxpool.Pool, logger *slog.Logger) *UsersStorage {
	return &UsersStorage{
		queries: db.New(pool),
		logger:  logger,
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

		u.logger.Error("failed to get user from database", "error", err)
		return domain.User{}, err
	}

	return domain.User{
		Login:    dbUser.Login.String,
		Password: dbUser.Password.String,
		Role:     dbUser.Role.String,
	}, nil
}
