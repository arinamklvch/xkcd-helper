package usecase

import "github.com/arinamklvch/xkcd-helper/internal/domain"

func (s *Service) GetLastComic() (domain.Comic, error) {
	lastComic, err := s.comicsStorage.GetLastComic()
	if err != nil {
		return domain.Comic{}, err
	}

	return lastComic, nil
}
