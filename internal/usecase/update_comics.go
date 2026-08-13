package usecase

import (
	"fmt"
	"slices"
	"strings"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/arinamklvch/xkcd-helper/pkg/utils"
)

func (s *Service) UpdateComics() error {
	// lastDbComic.Num == 0 when database is empty
	lastDbComic, err := s.comicsStorage.GetLastComic()
	if err != nil {
		return err
	}

	lastXkcdNum, err := s.xkcdCLient.GetLastComicNum()
	if err != nil {
		return err
	}

	if lastXkcdNum <= lastDbComic.Num {
		fmt.Println("comics are up-to-date")
		return nil
	}

	comics, err := s.xkcdCLient.DownloadComicsRange(lastDbComic.Num+1, lastXkcdNum)
	if err != nil {
		return err
	}

	if err := s.comicsStorage.InsertComics(comics); err != nil {
		return err
	}

	return s.insertIntoInvertedIndex(comics)
}

func (s *Service) insertIntoInvertedIndex(comics []domain.Comic) error {
	// building inverted index
	wordsToNums := make(map[string][]int32)
	for _, comic := range comics {
		text := strings.Join([]string{comic.SafeTitle, comic.Transcript, comic.Alt}, " ")
		words, err := utils.NormalizeWords(text)
		if err != nil {
			return err
		}

		for _, word := range words {
			if !slices.Contains(wordsToNums[word], int32(comic.Num)) {
				wordsToNums[word] = append(wordsToNums[word], int32(comic.Num))
			}
		}
	}

	args := make([]adapter.InsertIntoInvertedIndexArgs, 0, len(wordsToNums))
	for word, nums := range wordsToNums {
		args = append(args, adapter.InsertIntoInvertedIndexArgs{
			Word:       word,
			ComicsNums: nums,
		})
	}

	return s.invertedIndexStorage.InsertIntoInvertedIndex(args)
}
