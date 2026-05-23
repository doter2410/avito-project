package courier

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreateCourier(w http.ResponseWriter, r *http.Request) {
	var input Courier
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.Phone == "" {
		http.Error(w, `{"error": "name and phone ere required"}`, http.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateCourier(r.Context(), input)
	if err != nil {
		http.Error(w, `{"error":"failed to create courier"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]int64{
		"id": id,
	})
}

func (h *Handler) GetCourier(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id format"}`, http.StatusBadRequest)
		return
	}

	c, err := h.storage.GetCourierById(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"courier not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) GetAllCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.storage.GetCouriers(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to get couriers"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (h *Handler) PutUpdCourier(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id format"}`, http.StatusBadRequest)
		return
	}

	var updCourier Courier
	if err := json.NewDecoder(r.Body).Decode(&updCourier); err != nil {
		http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
		return
	}

	err = h.storage.UpdateCourier(r.Context(), id, updCourier)
	if err != nil {
		http.Error(w, `{"error":"courier not found or update failed"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
