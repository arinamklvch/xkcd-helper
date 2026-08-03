package adapter

import (
	"context"

	"github.com/arinamklvch/xkcd-helper/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvertedIndexStorage struct {
	queries *db.Queries
}

type InsertIntoInvertedIndexArgs struct {
	Word       string
	ComicsNums []int32
}

func NewInvertedIndexStorage(pool *pgxpool.Pool) *InvertedIndexStorage {
	return &InvertedIndexStorage{
		queries: db.New(pool),
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

	return err
}

func (i *InvertedIndexStorage) GetFromInvertedIndex(words []string) ([][]int32, error) {
	comicsNums, err := i.queries.GetFromInvertedIndex(context.Background(), words)
	if err != nil {
		return nil, err
	}

	return comicsNums, nil
}
