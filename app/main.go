package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shatrovich/atnr.pro/service"
	"github.com/shatrovich/atnr.pro/storage"

	configuration "github.com/shatrovich/atnr.pro/config"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	config := configuration.NewDefault()

	logger, _ := zap.NewProduction()

	db, err := sql.Open("postgres", config.Database.GetAddr())

	if err != nil {
		logger.Fatal("error open database", zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		logger.Fatal("error ping database", zap.Error(err))
	}

	store := storage.NewStore(db, logger)

	service := service.NewService(store, logger)

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Timeout(time.Minute))

	router.Post("/", service.PushTaskHTTP)

	router.Get("/{id}", service.GetTaskDataHTTP)

	logger.Fatal("error listen http server", zap.Error(http.ListenAndServe(config.Application.GetAddr(), router)))
}
