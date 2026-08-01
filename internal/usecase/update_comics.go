package usecase

import (
	"slices"
	"strings"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

func (s *Service) UpdateComics() error {
	// latestDbNum == 0 when database is empty
	latestComic, err := s.comicsStorage.GetLatestComic()
	latestDbNum := latestComic.Num
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

	return s.insertIntoInvertedIndex(comics)
}

func (s *Service) insertIntoInvertedIndex(comics []domain.Comic) error {
	// building inverted index
	wordsToNums := make(map[string][]int32)
	for _, comic := range comics {
		allWords := strings.Split(comic.SafeTitle+" "+comic.Transcript+" "+comic.Alt, " ")
		for _, word := range allWords {
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
