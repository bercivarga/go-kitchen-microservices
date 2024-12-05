package main

import (
	"errors"
	common "github.com/bercivarga/commons"
	"log"
	"net/http"

	pb "github.com/bercivarga/commons/api"
)

type Handler struct {
	orderServiceClient pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *Handler {
	return &Handler{
		orderServiceClient: client,
	}
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("GET /health", h.HealthCheck)

	// Business logic
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.HandleCreateOrder)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

func (h *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.orderServiceClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Items:      items,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	common.WriteJSON(w, http.StatusCreated, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items must not be empty")
	}

	for _, item := range items {
		if item.Id == "" {
			return errors.New("item id must not be empty")
		}

		if item.Quantity < 1 {
			return errors.New("quantity must be greater than 0")
		}
	}

	return nil
}
