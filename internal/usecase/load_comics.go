package usecase

import (
	"errors"
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

var ErrInvalidRange = errors.New("invalid range")

func (s *Service) LoadComics(input dto.LoadComicsInput) ([]domain.Comic, error) {
	if input.From <= 0 || input.To <= 0 || input.From > input.To {
		return nil, ErrInvalidRange
	}

	comics, err := s.comicsStorage.GetComicsRange(input.From, input.To)
	if err != nil {
		return nil, fmt.Errorf("get comics from database: %w", err)
	}

	return comics, nil
}
