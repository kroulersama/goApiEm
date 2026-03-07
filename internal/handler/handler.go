package handler

import (
	"encoding/json"
	"errors"
	"goApiEM/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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
	// mux.HandleFunc("GET /subs", h.GetSubs)
	mux.HandleFunc("Get /subs/{id}/prices", h.GetPrices)
}

// CreateSubs - создаёт подписку
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

// GetSubByID - выдаёт подписку по id
func (h *handler) GetSubByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByIDSub(id)
	if err != nil {
		if errors.Is(err, nil) {
			http.Error(w, "подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// UpdareSub - обновляет действующую подписку
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
		if errors.Is(err, nil) {
			http.Error(w, "подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// DeleteSubs - удаляет подписку по id
func (h *handler) DeleteSubs(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	if err = h.service.DeleteSub(id); err != nil {
		if errors.Is(err, nil) {
			http.Error(w, "подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPrices - возврящает сумму подписки за период по id юзера и id подписки
func (h *handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "неккоректный id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name      string     `json:"name"`
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
		if errors.Is(err, nil) {
			http.Error(w, "подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(prices); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
