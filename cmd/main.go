// Package main provides an HTTP API for loading XKCD comics by range.
//
// @title XKCD Helper API
// @version 1.0
// @description API for loading XKCD comic titles by comic number range.
package main

import (
	"context"
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
)

func main() {
	xkcdClient := adapter.NewXkcdClient(*http.DefaultClient)
	service := usecase.New(xkcdClient)
	router := controller.NewRouter(service)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server failed:", err)
		}
	}()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	s := <-signalChannel
	fmt.Println("\nCatched signal:", s)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
