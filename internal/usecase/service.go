package usecase

import "github.com/arinamklvch/xkcd-helper/internal/adapter"

type Service struct {
	xkcdCLient           *adapter.XkcdClient
	comicsStorage        *adapter.ComicsStorage
	invertedIndexStorage *adapter.InvertedIndexStorage
	usersStorage         *adapter.UsersStorage
}

func New(xkcdClient *adapter.XkcdClient, comicsStorage *adapter.ComicsStorage, invertedIndexStorage *adapter.InvertedIndexStorage, usersStorage *adapter.UsersStorage) *Service {
	return &Service{
		xkcdCLient:           xkcdClient,
		comicsStorage:        comicsStorage,
		invertedIndexStorage: invertedIndexStorage,
		usersStorage:         usersStorage,
	}
}
