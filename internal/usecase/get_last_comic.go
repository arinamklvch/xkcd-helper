package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

func (s *Service) GetLastComic() (domain.Comic, error) {
	lastComic, err := s.comicsStorage.GetLastComic()
	if err != nil {
		return domain.Comic{}, fmt.Errorf("get last comic from database: %w", err)
	}

	return lastComic, nil
}
