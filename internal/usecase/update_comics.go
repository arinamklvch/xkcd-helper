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
		return fmt.Errorf("get last comic from database: %w", err)
	}

	lastXkcdNum, err := s.xkcdCLient.GetLastComicNum()
	if err != nil {
		return fmt.Errorf("get last comic from xkcd: %w", err)
	}

	if lastXkcdNum <= lastDbComic.Num {
		s.logger.Info("comics are up-to-date", "last_comic_num", lastDbComic.Num)
		return nil
	}

	comics, err := s.xkcdCLient.DownloadComicsRange(lastDbComic.Num+1, lastXkcdNum)
	if err != nil {
		return fmt.Errorf("download comics from xkcd: %w", err)
	}

	if err := s.comicsStorage.InsertComics(comics); err != nil {
		return fmt.Errorf("insert comics into database: %w", err)
	}

	if err := s.insertIntoInvertedIndex(comics); err != nil {
		return fmt.Errorf("insert comics into inverted index: %w", err)
	}

	return nil
}

func (s *Service) insertIntoInvertedIndex(comics []domain.Comic) error {
	// building inverted index
	wordsToNums := make(map[string][]int32)
	for _, comic := range comics {
		text := strings.Join([]string{comic.SafeTitle, comic.Transcript, comic.Alt}, " ")
		words, err := utils.NormalizeWords(text)
		if err != nil {
			return fmt.Errorf("normalize comic words: %w", err)
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

	if err := s.invertedIndexStorage.InsertIntoInvertedIndex(args); err != nil {
		return fmt.Errorf("insert inverted index into database: %w", err)
	}

	return nil
}
