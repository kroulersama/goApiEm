package config

import (
	"log"
	"net/http"
	"time"
)

// responseWriter обёртка для http.ResponseWriter чтобы перехватить статус код
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает запись статус кода
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware логирует все входящие HTTP запросы
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем входящий запрос
		log.Printf("[ВХОД] %s %s", r.Method, r.URL.Path)

		// Создаём обёртку для ResponseWriter
		wrapper := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Вызываем следующий обработчик
		next.ServeHTTP(wrapper, r)

		// Логируем результат
		duration := time.Since(start)
		log.Printf("[ВЫХОД] %s %s - статус: %d, время: %v",
			r.Method, r.URL.Path, wrapper.statusCode, duration)
	})
}

// InitLogger инициализирует базовые настройки логирования
func InitLogger() {
	// Можно добавить настройки формата логов
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("✅ Логгер инициализирован")
}
