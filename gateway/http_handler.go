package main

import "net/http"

type handler struct {
	// gateway
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(writer http.ResponseWriter, request *http.Request) {

}
