package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/manoelduran/go-with-react.git/internal/api"
	"github.com/manoelduran/go-with-react.git/internal/store/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	ctx := context.Background()
	pool, err := pgx.Connect(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",os.Getenv("WSRS_DB_USER"), os.Getenv("WSRS_DB_PASSWORD"), os.Getenv("WSRS_DB_HOST"), os.Getenv("WSRS_DB_PORT"), os.Getenv("WSRS_DB_NAME")))
	if err != nil {
		panic(err)
	}
	defer pool.Close(ctx)

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	handler := api.NewHandler(pgstore.New(pool))
	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()
	quit := make(chan os.Signal,1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
