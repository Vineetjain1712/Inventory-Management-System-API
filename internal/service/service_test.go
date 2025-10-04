package service

import (
	"os"
	"testing"

	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/models"
	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/store"
)

var svc *Service

func TestMain(m *testing.M) {
	// Initialize a temporary SQLite DB for tests
	db, err := store.NewStore("test_inventory.db")
	if err != nil {
		panic(err)
	}

	// Clear table before running tests
	db.DB.Exec("DELETE FROM products")

	svc = NewService(db)

	// Run tests
	code := m.Run()

	// Clean up DB file
	os.Remove("test_inventory.db")
	os.Exit(code)
}

func TestCreateProduct(t *testing.T) {
	// Valid product
	p := &models.Product{
		Name:              "Test Product",
		Description:       "A test product",
		StockQuantity:     5,
		LowStockThreshold: 2,
	}
	id, err := svc.CreateProduct(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatal("expected valid ID, got 0")
	}

	// Invalid product (negative stock)
	p2 := &models.Product{
		Name:          "Bad Product",
		StockQuantity: -5,
	}
	_, err = svc.CreateProduct(p2)
	if err == nil {
		t.Fatal("expected error for negative stock, got nil")
	}
}

func TestIncreaseStock(t *testing.T) {
	// Create product
	p := &models.Product{
		Name:              "StockProduct",
		StockQuantity:     5,
		LowStockThreshold: 2,
	}
	id, _ := svc.CreateProduct(p)

	// Increase stock valid
	updated, err := svc.IncreaseStock(id, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.StockQuantity != 8 {
		t.Fatalf("expected stock 8, got %d", updated.StockQuantity)
	}

	// Increase stock invalid (negative amount)
	_, err = svc.IncreaseStock(id, -3)
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
}

func TestDecreaseStock(t *testing.T) {
	// Create product
	p := &models.Product{
		Name:              "StockProduct2",
		StockQuantity:     5,
		LowStockThreshold: 2,
	}
	id, _ := svc.CreateProduct(p)

	// Decrease stock valid
	updated, err := svc.DecreaseStock(id, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.StockQuantity != 2 {
		t.Fatalf("expected stock 2, got %d", updated.StockQuantity)
	}

	// Decrease stock invalid (more than available)
	_, err = svc.DecreaseStock(id, 5)
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}

	// Decrease stock invalid (negative amount)
	_, err = svc.DecreaseStock(id, -2)
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
}

func TestListLowStockProducts(t *testing.T) {
	// Create products
	p1 := &models.Product{Name: "P1", StockQuantity: 1, LowStockThreshold: 2}
	p2 := &models.Product{Name: "P2", StockQuantity: 3, LowStockThreshold: 2}
	svc.CreateProduct(p1)
	svc.CreateProduct(p2)

	lowStock, err := svc.ListLowStockProducts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(lowStock) != 1 || lowStock[0].Name != "P1" {
		t.Fatalf("expected 1 low stock product P1, got %v", lowStock)
	}
}
