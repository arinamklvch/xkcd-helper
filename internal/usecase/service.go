package usecase

import "github.com/arinamklvch/xkcd-helper/internal/adapter"

type Service struct {
	xkcdCLient           *adapter.XkcdClient
	comicsStorage        *adapter.ComicsStorage
	invertedIndexStorage *adapter.InvertedIndexStorage
}

func New(xkcdClient *adapter.XkcdClient, comicsStorage *adapter.ComicsStorage, invertedIndexStorage *adapter.InvertedIndexStorage) *Service {
	return &Service{
		xkcdCLient:           xkcdClient,
		comicsStorage:        comicsStorage,
		invertedIndexStorage: invertedIndexStorage,
	}
}
