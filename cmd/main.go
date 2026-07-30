// Package main provides an HTTP API for loading and searching XKCD comics.
//
// @title XKCD Helper API
// @version 1.0
// @description API for loading XKCD comics by number range and searching comics by query words.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/controller"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("server stopped:", err)
		os.Exit(1)
	}
}

func run() error {
	// xkcdClient -- штука которая идет в xkcd и скачивает комиксы
	// создаем XkcdClient (с дефолтными настройками)
	xkcdClient := adapter.NewXkcdClient(*http.DefaultClient)

	// подключение к базе
	databaseURL := "postgres://admin:admin@localhost:5433/xkcd?sslmode=disable"
	pool, err := initPostgreSQL(databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// comicsStorage -- штука которая умеет удобно отправлять запросы в БД
	// создаем ComicsStorage на основе подключения pool
	comicsStorage := adapter.NewComicsStorage(pool)
	invertedIndexStorage := adapter.NewInvertedIndexStorage(pool)

	// Service -- объект, в котором будут методы use case / бизнес-логики
	// xkcdClient, comicsStorage -- инструменты для похода во внешние источники
	service := usecase.New(xkcdClient, comicsStorage, invertedIndexStorage)
	// загружаем все/новые комиксы один раз при запуске
	err = service.UpdateComics()
	if err != nil {
		return err
	}

	// создаем HTTP-роутер
	// передаем в него service, чтобы хендлеры могли вызывать бизнес-логику
	router := controller.NewRouter(service)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		// т.к. код останавливается пока сервер работает без остановки и ошибок
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	// ждем либо сигнал, либо ошибку сервера
	select {
	case err := <-serverErrCh:
		if err != nil {
			return err
		}
	case s := <-signalCh:
		fmt.Println("\ncatched signal:", s)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func runMigrations(databaseURL string) error {
	// db -- объект доступа к БД
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, "migrations")
}

func initPostgreSQL(databaseURL string) (*pgxpool.Pool, error) {
	// ограничение времени на создание/инициализацию пула + проверку БД
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	pool, err := pgxpool.New(dbCtx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(dbCtx); err != nil {
		return nil, err
	}

	// применяем миграции
	if err := runMigrations(databaseURL); err != nil {
		return nil, err
	}

	return pool, nil
}
