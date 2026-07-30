package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

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
