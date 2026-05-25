package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/doter2410/avito-project/internal/model"
	"github.com/go-chi/chi/v5"
)

type CourierService interface {
	CreateCourier(ctx context.Context, c model.Courier) (int64, error)
	GetCourierById(ctx context.Context, id int64) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]*model.Courier, error)
	UpdateCourier(ctx context.Context, id int64, c model.Courier) error

	AssignDelivery(ctx context.Context, orderID string) (*model.Courier, *model.Delivery, error)
	UnassignDelivery(ctx context.Context, orderID string) (*model.Courier, error)
}

type CourierHandler struct {
	service CourierService
}

func NewCourierHandler(service CourierService) *CourierHandler {
	return &CourierHandler{service: service}
}

type DeliveryRequest struct {
	OrderID string `json:"order_id"`
}

type AssignResponse struct {
	CourierID        int64     `json:"courier_id"`
	OrderID          string    `json:"order_id"`
	TransportType    string    `json:"transport_type"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}

type UnassignResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int64  `json:"courier_id"`
}

func (h *CourierHandler) CreateCourier(w http.ResponseWriter, r *http.Request) {
	var input model.Courier
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.Phone == "" {
		http.Error(w, `{"error": "name and phone ere required"}`, http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateCourier(r.Context(), input)
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

func (h *CourierHandler) GetCourier(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id format"}`, http.StatusBadRequest)
		return
	}

	c, err := h.service.GetCourierById(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"courier not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *CourierHandler) GetAllCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.service.GetCouriers(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to get couriers"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (h *CourierHandler) PutUpdCourier(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id format"}`, http.StatusBadRequest)
		return
	}

	var updCourier model.Courier
	if err := json.NewDecoder(r.Body).Decode(&updCourier); err != nil {
		http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
		return
	}

	err = h.service.UpdateCourier(r.Context(), id, updCourier)
	if err != nil {
		http.Error(w, `{"error":"courier not found or update failed"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CourierHandler) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	var req DeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	c, d, err := h.service.AssignDelivery(r.Context(), req.OrderID)
	if err != nil {
		fmt.Println("РЕАЛЬНАЯ ОШИБКА БАЗЫ:", err)
		http.Error(w, "conflict or error", http.StatusConflict)
		return
	}
	resp := AssignResponse{
		CourierID:        c.ID,
		OrderID:          d.OrderID,
		TransportType:    c.TransportType,
		DeliveryDeadline: d.Deadline,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CourierHandler) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	var req DeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	c, err := h.service.UnassignDelivery(r.Context(), req.OrderID)
	if err != nil {
		http.Error(w, "conflict or error", http.StatusNotFound)
		return
	}
	resp := UnassignResponse{
		OrderID:   req.OrderID,
		Status:    "unassigned",
		CourierID: c.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
