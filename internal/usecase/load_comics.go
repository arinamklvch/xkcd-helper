package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

type Service struct {
	xkcdCLient    *adapter.XkcdClient
	comicsStorage *adapter.ComicsStorage
}

func New(xkcdClient *adapter.XkcdClient, comicsStorage *adapter.ComicsStorage) *Service {
	return &Service{
		xkcdCLient:    xkcdClient,
		comicsStorage: comicsStorage,
	}
}

func (s *Service) LoadComics(input dto.LoadComicsInput) ([]domain.Comic, error) {
	totalCnt := input.To - input.From + 1
	if totalCnt <= 0 {
		return nil, fmt.Errorf("invalid range")
	}

	comics, err := s.comicsStorage.GetComicsRange(input.From, input.To)
	if err != nil {
		return nil, fmt.Errorf("error getting comics from database")
	}

	return comics, nil
}

func (s *Service) UpdateComics() error {
	// latestDbNum == 0 when database is empty
	latestDbNum, err := s.comicsStorage.GetLatestComicNum()
	if err != nil {
		return err
	}

	latestNum, err := s.xkcdCLient.GetLatestComicNum()
	if err != nil {
		return err
	}

	comics, err := s.xkcdCLient.DownloadComicsRange(latestDbNum+1, latestNum)
	if err != nil {
		return err
	}

	// inserting into database
	if err := s.comicsStorage.InsertComics(comics); err != nil {
		return err
	}

	return nil
}
