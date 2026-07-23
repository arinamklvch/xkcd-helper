package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

type Service struct {
	xkcd *adapter.XkcdClient
}

func New(xkcdClient *adapter.XkcdClient) *Service {
	return &Service{
		xkcd: xkcdClient,
	}
}

func (s *Service) LoadComics(input dto.LoadComicsInput) ([]domain.Comic, error) {
	totalCnt := input.To - input.From + 1
	if totalCnt <= 0 {
		return []domain.Comic{}, fmt.Errorf("Invalid range")
	}

	comics, err := s.xkcd.DownloadComics(input.From, input.To)
	if err != nil {
		return []domain.Comic{}, fmt.Errorf("Error getting comics")
	}

	return comics, nil
}
