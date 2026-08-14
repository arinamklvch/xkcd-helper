package usecase

import (
	"fmt"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/arinamklvch/xkcd-helper/pkg/utils"
)

func (s *Service) SearchComics(input dto.SearchComicsInput) ([]domain.Comic, error) {
	queryWords, err := utils.NormalizeWords(input.Query)
	if err != nil {
		return nil, err
	}
	comicsNums, err := s.invertedIndexStorage.GetFromInvertedIndex(queryWords)
	if err != nil {
		return nil, fmt.Errorf("error while searching comics in database")
	}

	seen := make(map[int32]int)
	suitableNums := make([]int32, 0, s.maxFoundComics)
	length := len(queryWords)

OuterLoop:
	for _, arr := range comicsNums {
		for _, num := range arr {
			// пересечение всех query words
			if seen[num] >= length-1 {
				suitableNums = append(suitableNums, num)
				if len(suitableNums) == s.maxFoundComics {
					break OuterLoop
				}
			}
			seen[num]++
		}
	}

	// дополняем до maxSearchedComics
	for _, arr := range comicsNums {
		for i := 0; i < len(arr) && len(suitableNums) != s.maxFoundComics; i++ {
			num := arr[i]
			if seen[num] < length {
				suitableNums = append(suitableNums, num)
			}
		}
	}

	comics, err := s.comicsStorage.GetComicsByNums(suitableNums)
	if err != nil {
		return nil, err
	}

	return comics, nil
}
