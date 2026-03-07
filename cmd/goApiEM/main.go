package main

import (
	"goApiEM/internal/config"
	"goApiEM/internal/handler"
	"goApiEM/internal/repository"
	"goApiEM/internal/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Printf("Файл .env не найден")
	}

	config.InitLogger()

	dbConfig := &repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := repository.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	repo := repository.NewSubRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	service := service.NewSubSevrice(repo)
	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	loggedHandler := config.LoggerMiddleware(mux)
	log.Printf("Сервер запущен на http://localhost:%s", serverPort)

	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      loggedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
