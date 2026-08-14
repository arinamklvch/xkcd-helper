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
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/config"
	"github.com/arinamklvch/xkcd-helper/internal/controller"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/time/rate"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("server stopped:", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.New()
	if err != nil {
		return err
	}

	// определение уровня логирования
	logLevel, err := parseLogLevel(config.LoggerLevel)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	xkcdClient := adapter.NewXkcdClient(*http.DefaultClient, config.MaxWorkers, logger)

	// подключение к базе
	pool, err := initPostgreSQL(config.DatabaseURL, config.DbTimeout)
	if err != nil {
		return err
	}
	defer pool.Close()

	// создаем storages на основе подключения pool
	comicsStorage := adapter.NewComicsStorage(pool, logger)
	invertedIndexStorage := adapter.NewInvertedIndexStorage(pool, logger)
	usersStorage := adapter.NewUsersStorage(pool, logger)

	// Service -- объект, в котором будут методы use case / бизнес-логики
	// + инструменты для похода во внешние источники
	service := usecase.New(xkcdClient, comicsStorage, invertedIndexStorage, usersStorage,
		config.TokenTTL, config.MaxFoundComics, config.JWTsecretKey, logger)

	// загружаем все/новые комиксы один раз при запуске
	err = service.UpdateComics()
	if err != nil {
		return err
	}

	// создаем HTTP-роутер
	// передаем в него service, чтобы хендлеры могли вызывать бизнес-логику
	router := controller.NewRouter(service, rate.Limit(config.RateLimit), config.RateBurst, config.JWTsecretKey, logger)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(config.ServerTimeout)*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func parseLogLevel(value string) (slog.Level, error) {
	var logLevel slog.Level

	if value == "" {
		return slog.LevelInfo, nil
	}

	err := logLevel.UnmarshalText([]byte(strings.ToUpper(value)))
	if err != nil {
		return slog.LevelInfo, fmt.Errorf("invalid logger_level %q: %w", value, err)
	}

	return logLevel, nil
}

func runMigrations(databaseURL string) error {
	// db -- объект доступа к БД
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Println("failed to close db:", closeErr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, "migrations")
}

func initPostgreSQL(databaseURL string, timeout int) (*pgxpool.Pool, error) {
	// ограничение времени на создание/инициализацию пула + проверку БД
	dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
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
