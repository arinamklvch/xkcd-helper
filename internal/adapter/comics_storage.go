package adapter

import (
	"context"
	"errors"

	"github.com/arinamklvch/xkcd-helper/internal/db"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ComicsStorage struct {
	queries *db.Queries
}

func NewComicsStorage(pool *pgxpool.Pool) *ComicsStorage {
	return &ComicsStorage{
		queries: db.New(pool),
	}
}

func (c *ComicsStorage) GetLastComic() (domain.Comic, error) {
	lastComic, err := c.queries.GetLastComic(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Comic{}, nil
		}
		return domain.Comic{}, err
	}

	return mapDbComicToDomainComic(lastComic), nil
}

func (c *ComicsStorage) GetComicsRange(from, to int) ([]domain.Comic, error) {
	dbComics, err := c.queries.GetComicsRange(context.Background(), db.GetComicsRangeParams{
		From: pgtype.Int4{Int32: int32(from), Valid: true},
		To:   pgtype.Int4{Int32: int32(to), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return mapDbComicsToDomainComics(dbComics), nil
}

func (c *ComicsStorage) GetComicsByNums(nums []int32) ([]domain.Comic, error) {
	dbComics, err := c.queries.GetComicsByNums(context.Background(), nums)
	if err != nil {
		return nil, err
	}

	return mapDbComicsToDomainComics(dbComics), nil
}

func mapDbComicToDomainComic(dbComic db.Comic) domain.Comic {
	return domain.Comic{
		Month:      dbComic.Month.String,
		Num:        int(dbComic.Num.Int32),
		Link:       dbComic.Link.String,
		Year:       dbComic.Year.String,
		News:       dbComic.News.String,
		SafeTitle:  dbComic.SafeTitle.String,
		Transcript: dbComic.Transcript.String,
		Alt:        dbComic.Alt.String,
		Img:        dbComic.Img.String,
		Title:      dbComic.SafeTitle.String,
		Day:        dbComic.Day.String,
	}
}

func mapDbComicsToDomainComics(dbComics []db.Comic) []domain.Comic {
	domainComics := make([]domain.Comic, 0, len(dbComics))
	for _, comic := range dbComics {
		domainComics = append(domainComics, mapDbComicToDomainComic(comic))
	}
	return domainComics
}

func (c *ComicsStorage) InsertComics(comics []domain.Comic) error {
	dbComics := make([]db.InsertComicsParams, 0, len(comics))
	for _, comic := range comics {
		dbComics = append(dbComics, db.InsertComicsParams{
			Month:      pgtype.Text{String: comic.Month, Valid: true},
			Num:        pgtype.Int4{Int32: int32(comic.Num), Valid: true},
			Link:       pgtype.Text{String: comic.Link, Valid: true},
			Year:       pgtype.Text{String: comic.Year, Valid: true},
			News:       pgtype.Text{String: comic.News, Valid: true},
			SafeTitle:  pgtype.Text{String: comic.SafeTitle, Valid: true},
			Transcript: pgtype.Text{String: comic.Transcript, Valid: true},
			Alt:        pgtype.Text{String: comic.Alt, Valid: true},
			Img:        pgtype.Text{String: comic.Img, Valid: true},
			Title:      pgtype.Text{String: comic.Title, Valid: true},
			Day:        pgtype.Text{String: comic.Day, Valid: true},
		})
	}

	_, err := c.queries.InsertComics(context.Background(), dbComics)

	return err
}
