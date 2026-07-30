package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

const maxSearchedComics = 10

func (s *Service) SearchComics(input dto.SearchComicsInput) ([]domain.Comic, error) {
	comicsNums, err := s.invertedIndexStorage.GetFromInvertedIndex(input.Words)
	if err != nil {
		return nil, fmt.Errorf("error while searching comics in database")
	}

	seen := make(map[int32]int)
	suitableNums := make([]int32, 0, maxSearchedComics)
	length := len(input.Words)
	// найти пересечение, если там < maxSearchedComics то добить рандомными
OuterLoop:
	for _, arr := range comicsNums {
		for _, num := range arr {
			if seen[num] >= length-1 {
				suitableNums = append(suitableNums, num)
				if len(suitableNums) == maxSearchedComics {
					break OuterLoop
				}
			}
			seen[num]++
		}
	}

	for _, arr := range comicsNums {
		for _, num := range arr {
			if seen[num] < length-1 {
				suitableNums = append(suitableNums, num)
			}
		}
		if len(suitableNums) >= maxSearchedComics {
			suitableNums = suitableNums[:maxSearchedComics]
			break
		}
	}

	comics, err := s.comicsStorage.GetComicsByNums(suitableNums)
	if err != nil {
		return nil, err
	}

	return comics, nil
}
