package usecase

import "github.com/arinamklvch/xkcd-helper/internal/domain"

func (s *Service) GetLatestComic() (domain.Comic, error) {
	latestComic, err := s.comicsStorage.GetLatestComic()
	if err != nil {
		return domain.Comic{}, err
	}

	return latestComic, nil
}
