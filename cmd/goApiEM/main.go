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
	"gorm.io/gorm/logger"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Printf("Файл .env не найден")
	}

	// Инициализация логгера
	config.InitLogger()
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	dbConfig := &repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Подеключение к базе
	db, err := repository.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	db.Logger = dbLogger

	// Создание репозитория
	repo := repository.NewSubRepo(db)
	// Подключение миграций
	if err := repo.AutoMigrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// Создание сервиса и хэндлера
	service := service.NewSubSevrice(repo)
	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	// Подлючение путей
	handler.RegisterRoutes(mux)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Логирование
	loggedHandler := config.LoggerMiddleware(mux)
	log.Printf("Сервер запущен на http://localhost:%s", serverPort)

	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      loggedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Чтение запросов
	log.Fatal(server.ListenAndServe())

}
