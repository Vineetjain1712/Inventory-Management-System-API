// Package service
package service

import (
	"database/sql"
	"errors"
	"math"
	"strings"

	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/models"
	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/store"
)

// Public, reusable errors (map these to HTTP codes in your handlers)
var (
	ErrNilProduct         = errors.New("product is nil")              // 400
	ErrInvalidID          = errors.New("invalid id")                  // 400
	ErrEmptyName          = errors.New("name is required")            // 400
	ErrNegativeStock      = errors.New("stock_quantity cannot be negative")      // 400
	ErrNegativeThreshold  = errors.New("low_stock_threshold cannot be negative") // 400
	ErrZeroOrNegativeAmt  = errors.New("amount must be greater than zero")       // 400
	ErrInsufficientStock  = errors.New("insufficient stock")          // 400
	ErrStockOverflow      = errors.New("stock quantity overflow")     // 400
	ErrNotFound           = errors.New("product not found")           // 404
)

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

// ----------------- CRUD Operations -----------------


// CreateProduct inserts a new product row and returns the auto-generated ID.
func (s *Service) CreateProduct(p *models.Product) (int64, error) {
	if err := validateProductForCreate(p); err != nil {
		return 0, err
	}
	return s.store.CreateProduct(p)
}

// GetProduct returns a single product.
func (s *Service) GetProduct(id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	p, err := s.store.GetProduct(id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return p, nil
}

// ListProducts returns all products.
func (s *Service) ListProducts() ([]models.Product, error) {
	return s.store.ListProducts()
}


// UpdateProduct validates the payload then updates the row.
// Returns ErrNotFound (404) if the product does not exist.
func (s *Service) UpdateProduct(p *models.Product) error {
	if p == nil {
		return ErrNilProduct
	}
	if p.ID <= 0 {
		return ErrInvalidID
	}
	if err := validateCommonProductFields(p); err != nil {
		return err
	}
	ok, err := s.store.UpdateProduct(p)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// DeleteProduct removes the product if it exists.
// Returns ErrNotFound (404) if no such product.
func (s *Service) DeleteProduct(id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}
	ok, err := s.store.DeleteProduct(id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// ----------------- Inventory Operations -----------------

func (s *Service) IncreaseStock(id int64, amount int) (*models.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	if amount <= 0 {
		return nil, ErrZeroOrNegativeAmt
	}

	p, err := s.store.GetProduct(id)
	if err != nil {
		return nil, mapNotFound(err)
	}

	// Overflow guard
	if amount > 0 && p.StockQuantity > (math.MaxInt-amount) {
		return nil, ErrStockOverflow
	}

	p.StockQuantity += amount
	if _,err := s.store.UpdateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DecreaseStock(id int64, amount int) (*models.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	if amount <= 0 {
		return nil, ErrZeroOrNegativeAmt
	}

	p, err := s.store.GetProduct(id)
	if err != nil {
		return nil, mapNotFound(err)
	}

	if p.StockQuantity < amount {
		return nil, ErrInsufficientStock
	}

	p.StockQuantity -= amount
	if _,err := s.store.UpdateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

// ----------------- Bonus: Low Stock -----------------

func (s *Service) ListLowStockProducts() ([]models.Product, error) {
	products, err := s.store.ListProducts()
	if err != nil {
		return nil, err
	}
	// threshold <= 0 means "no threshold"
	low := make([]models.Product, 0, len(products))
	for _, p := range products {
		if p.LowStockThreshold > 0 && p.StockQuantity < p.LowStockThreshold {
			low = append(low, p)
		}
	}
	return low, nil
}

// ----------------- Helpers -----------------

func validateProductForCreate(p *models.Product) error {
	if p == nil {
		return ErrNilProduct
	}
	return validateCommonProductFields(p)
}

func validateCommonProductFields(p *models.Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrEmptyName
	}
	if p.StockQuantity < 0 {
		return ErrNegativeStock
	}
	if p.LowStockThreshold < 0 {
		return ErrNegativeThreshold
	}
	return nil
}

// Convert store-layer "not found" to service ErrNotFound
func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
