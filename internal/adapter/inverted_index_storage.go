package adapter

import (
	"context"
	"log/slog"

	"github.com/arinamklvch/xkcd-helper/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvertedIndexStorage struct {
	queries *db.Queries
	logger  *slog.Logger
}

type InsertIntoInvertedIndexArgs struct {
	Word       string
	ComicsNums []int32
}

func NewInvertedIndexStorage(pool *pgxpool.Pool, logger *slog.Logger) *InvertedIndexStorage {
	return &InvertedIndexStorage{
		queries: db.New(pool),
		logger:  logger,
	}
}

func (i *InvertedIndexStorage) InsertIntoInvertedIndex(args []InsertIntoInvertedIndexArgs) error {
	dbArgs := make([]db.InsertIntoInvertedIndexParams, 0, len(args))
	for _, arg := range args {
		dbArgs = append(dbArgs, db.InsertIntoInvertedIndexParams{
			Word:       pgtype.Text{String: arg.Word, Valid: true},
			ComicsNums: arg.ComicsNums,
		})
	}
	_, err := i.queries.InsertIntoInvertedIndex(context.Background(), dbArgs)
	if err != nil {
		i.logger.Error("failed to insert inverted index into database",
			"error", err,
		)
		return err
	}

	return nil
}

func (i *InvertedIndexStorage) GetFromInvertedIndex(words []string) ([][]int32, error) {
	comicsNums, err := i.queries.GetFromInvertedIndex(context.Background(), words)
	if err != nil {
		i.logger.Error("failed to get comics nums from inverted index",
			"error", err,
		)
		return nil, err
	}

	return comicsNums, nil
}
