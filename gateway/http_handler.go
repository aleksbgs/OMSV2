package main

import (
	"errors"
	"github.com/akeksbgs/omsv2-gateway/gateway"
	common "github.com/aleksbgs/commons"
	pb "github.com/aleksbgs/commons/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
}

func (h *handler) handleCreateOrder(writer http.ResponseWriter, r *http.Request) {

	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJson(r, &items); err != nil {
		common.WriteError(writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(writer, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	//grpc handling errors
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(writer, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJson(writer, http.StatusOK, o)
}
func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("item must have an ID")
		}
		if item.Quantity < 0 {
			return errors.New("item must have a non-negative quantity")
		}
	}

	return nil
}
