package handler

import (
	"encoding/json"
	"errors"
	"goApiEM/internal/repository"
	"goApiEM/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "goApiEM/docs"

	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// handler - HTTP-обработчик
type handler struct {
	service *service.SubService
}

// NewHandler - создаёт новый обработчик
func NewHandler(service *service.SubService) *handler { return &handler{service: service} }

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /subs", h.CreateSubs)
	mux.HandleFunc("GET /subs/{id}", h.GetSubByID)
	mux.HandleFunc("POST /subs/{id}", h.UpdateSub)
	mux.HandleFunc("DELETE /subs/{id}", h.DeleteSubs)
	mux.HandleFunc("GET /subs", h.GetSubs)
	mux.HandleFunc("GET /subs/{id}/prices", h.GetPrices)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

// CreateSubs
// @Summary      Создать подписку
// @Description  Создаёт новую подписку
// @Tags         subs
// @Accept       json
// @Produce      json
// @Param        request body object{name=string,price=int,user_id=string,start_date=string,end_date=string} true "Данные подписки"
// @Success      201  {object}  repository.Sub
// @Failure      400  {string}  string "некорректный JSON"
// @Failure      500  {string}  string "внутренняя ошибка"
// @Router       /subs [post]
func (h *handler) CreateSubs(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string     `json:"name"`
		Price     int        `json:"price"`
		UserID    uuid.UUID  `json:"user_id"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	sub, err := h.service.CreateSub(input.Name, input.Price,
		input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// GetSubByID
// @Summary      Получить подписку по ID
// @Description  Возвращает подписку по её идентификатору
// @Tags         subs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  repository.Sub
// @Failure      400  {string}  string "некорректный id"
// @Failure      404  {string}  string "подписка не найдена"
// @Failure      500  {string}  string "внутренняя ошибка"
// @Router       /subs/{id} [get]
func (h *handler) GetSubByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByIDSub(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, "подписка не найдена", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// UpdateSub
// @Summary      Обновить подписку
// @Description  Обновляет существующую подписку
// @Tags         subs
// @Accept       json
// @Produce      json
// @Param        id      path      int                                  true  "ID подписки"
// @Param        request body      object{service_name=string,price=int,user_id=string,start_date=string,end_date=string} true "Данные для обновления"
// @Success      200     {object}  repository.Sub
// @Failure      400     {string}  string "некорректный JSON или id"
// @Failure      404     {string}  string "подписка не найдена"
// @Failure      500     {string}  string "внутренняя ошибка"
// @Router       /subs/{id} [post]
func (h *handler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string     `json:"service_name"`
		Price     int        `json:"price"`
		UserID    uuid.UUID  `json:"user_id"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.UpdateSub(id, input.Name, input.Price, input.UserID,
		input.StartDate, input.EndDate)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, "подписка не найдена", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// DeleteSubs
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      204  "No Content"
// @Failure      400  {string}  string "некорректный id"
// @Failure      404  {string}  string "подписка не найдена"
// @Failure      500  {string}  string "внутренняя ошибка"
// @Router       /subs/{id} [delete]
func (h *handler) DeleteSubs(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	if err = h.service.DeleteSub(id); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, "подписка не найдена", http.StatusNotFound)
		case errors.Is(err, repository.ErrReqDB):
			http.Error(w, "Ошибка обращения к базе данных", http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSubs
// @Summary      Получить все подписки
// @Description  Возвращает список всех подписок
// @Tags         subs
// @Accept       json
// @Produce      json
// @Success      200  {array}   repository.Sub
// @Failure      500  {string}  string "внутренняя ошибка"
// @Router       /subs [get]
func (h *handler) GetSubs(w http.ResponseWriter, r *http.Request) {

	subs, err := h.service.GetSubs()
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReqDB):
			http.Error(w, "Ошибка обращения к базе данных", http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// GetPrices
// @Summary      Получить сумму за период
// @Description  Возвращает сумму всех подписок за указанный период
// @Tags         subs
// @Accept       json
// @Produce      json
// @Param        id           path      int    true  "ID подписки"
// @Param        request      body      object{user_id=string,start_date=string,end_date=string} true "Параметры периода"
// @Success      200          {object}  int64 "сумма"
// @Failure      400          {string}  string "некорректный id или JSON"
// @Failure      500          {string}  string "внутренняя ошибка"
// @Router       /subs/{id}/prices [get]
func (h *handler) GetPrices(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "неккоректный id", http.StatusBadRequest)
		return
	}

	var input struct {
		UserID    uuid.UUID  `json:"user_id"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	prices, err := h.service.GetPriceForRangeSub(id, input.UserID, input.StartDate, *input.EndDate)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReqDB):
			http.Error(w, "Ошибка обращения к базе данных", http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(prices); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
