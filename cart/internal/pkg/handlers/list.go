package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

// TODO добавить еще один меньший CartItem
type FullCartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type ListRequest struct {
	User int64
}

func (l ListRequest) Validate() error {
	if l.User < 0 {
		return ErrIncorrectUser
	}
	return nil
}

type ListResponse struct {
	Item       []FullCartItem
	TotalPrice uint32
}

type CartLister interface {
	GetList(ctx context.Context, user int64) (ListResponse, error)
}

type ListHandler struct {
	name       string
	cartLister CartLister
}

func NewListHandler(cartLister CartLister) *ListHandler {
	return &ListHandler{
		name:       "cart list handler",
		cartLister: cartLister,
	}
}

func (h ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &ListRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	resp, err := h.cartLister.GetList(r.Context(), req.User)
	if err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
}
