// Package handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/models"
	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/service"
	"github.com/gorilla/mux"
)

// ---- Error helpers ----

type apiError struct {
	Error string `json:"error"`
}

const (
	errInvalidBody   = "invalid request body"
	errInvalidID     = "invalid product ID"
	errInvalidAmount = "invalid amount"
	errNotFound      = "product not found"
)

func respondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, apiError{Error: msg})
}

// ProductHandler holds reference to service layer
type ProductHandler struct {
	Service *service.Service
}

// NewProductHandler returns a new ProductHandler
func NewProductHandler(svc *service.Service) *ProductHandler {
	return &ProductHandler{Service: svc}
}

// ---- small helper for ID parsing ----
func parseID(r *http.Request) (int64, error) {
	idParam := strings.TrimSpace(mux.Vars(r)["id"])
	return strconv.ParseInt(idParam, 10, 64)
}

// ----------------- CRUD Handlers -----------------

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	id, err := h.Service.CreateProduct(&p)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	p.ID = id
	respondJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	p, err := h.Service.GetProduct(id)
	if err != nil {
		respondError(w, http.StatusNotFound, errNotFound)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.ListProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list products")
		return
	}
	respondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	p.ID = id
	if err := h.Service.UpdateProduct(&p); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	if err := h.Service.DeleteProduct(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

// ----------------- Stock Handlers -----------------

func (h *ProductHandler) IncreaseStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	var req struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidAmount)
		return
	}

	p, err := h.Service.IncreaseStock(id, req.Amount)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	var req struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, errInvalidAmount)
		return
	}

	p, err := h.Service.DecreaseStock(id, req.Amount)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, p)
}

// ----------------- Bonus -----------------

func (h *ProductHandler) ListLowStockProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.ListLowStockProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list low stock products")
		return
	}
	respondJSON(w, http.StatusOK, products)
}
