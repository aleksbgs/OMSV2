package main

import (
	common "github.com/aleksbgs/commons"
	pb "github.com/aleksbgs/commons/api"
	"net/http"
)

type handler struct {
	// gateway
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
}

func (h *handler) handleCreateOrder(writer http.ResponseWriter, request *http.Request) {

	customerID := request.PathValue("customerID")

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJson(request, &items); err != nil {
		common.WriteError(writer, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.client.CreateOrder(request.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	if err != nil {
		common.WriteError(writer, http.StatusBadRequest, err.Error())
		return
	}
}
