package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type StocksService interface {
	GetStocks(sku uint32) uint64
}

type StocksHandler struct {
	stocksService StocksService
}

func NewStocksHandler(stocksService StocksService) *StocksHandler {
	return &StocksHandler{
		stocksService: stocksService,
	}
}

type StockRequest struct {
	SKU uint32 `json:"sku"`
}

var ErrIncorrectSKU = errors.New("incorrect SKU")

func (r StockRequest) Validate() error {
	if r.SKU == 0 {
		return ErrIncorrectSKU
	}
	return nil
}

type StockResponse struct {
	Count uint64 `json:"count"`
}

func (s StocksHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &StockRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Println("stocks: failed to decode request body")
		GetErrorResponse(w, "stocks", err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		log.Println("stocks: failed to validate request: %w", err)
		GetErrorResponse(w, "stocks", err, http.StatusBadRequest)
		return
	}

	count := s.stocksService.GetStocks(req.SKU)

	stockResponse := &StockResponse{
		Count: count,
	}
	raw, err := json.Marshal(stockResponse)
	if err != nil {
		GetErrorResponse(w, "stocks", err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	GetSuccessResponseWithBody(w, raw)
}
