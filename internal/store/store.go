// Package store
package store

import (
	"database/sql"
	"log"

	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/models"
	_ "modernc.org/sqlite"
)

// Store wraps the DB connection
type Store struct {
	DB *sql.DB
}

// NewStore initializes SQLite DB
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL CHECK (trim(name) <> ''),
		description TEXT,
		stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
		low_stock_threshold INTEGER NOT NULL DEFAULT 0 CHECK (low_stock_threshold >= 0)
		);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	log.Println("SQLite DB initialized at", dbPath)

	return &Store{DB: db}, nil
}

// CRUD Functions -
// CreateProduct inserts a new product row and returns the auto-generated ID.
// It relies on SQLite's INTEGER PRIMARY KEY to assign the id (no manual increments).
func (s *Store) CreateProduct(p *models.Product) (int64, error) {
	const query = `
		INSERT INTO products (name, description, stock_quantity, low_stock_threshold)
		VALUES (?, ?, ?, ?)`
	result, err := s.DB.Exec(query, p.Name, p.Description, p.StockQuantity, p.LowStockThreshold)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetProduct retrieves a product by its ID.
// Returns sql.ErrNoRows if no row exists, so callers can map to 404.
func (s *Store) GetProduct(id int64) (*models.Product, error) {
	const query = `
		SELECT id, name, description, stock_quantity, low_stock_threshold
		FROM products
		WHERE id = ?`
	row := s.DB.QueryRow(query, id)

	var p models.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StockQuantity, &p.LowStockThreshold); err != nil {
		return nil, err
	}
	return &p, nil
}

// ListProducts returns all products in the table.
// On large datasets you might add pagination later.
func (s *Store) ListProducts() ([]models.Product, error) {
	const query = `
		SELECT id, name, description, stock_quantity, low_stock_threshold
		FROM products`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.StockQuantity, &p.LowStockThreshold); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// UpdateProduct updates a product's fields by ID and reports whether any row was updated.
// Returns (false, nil) when the product does not exist (callers can map to 404).
func (s *Store) UpdateProduct(p *models.Product) (bool, error) {
	const query = `
		UPDATE products
		SET name = ?, description = ?, stock_quantity = ?, low_stock_threshold = ?
		WHERE id = ?`
	res, err := s.DB.Exec(query, p.Name, p.Description, p.StockQuantity, p.LowStockThreshold, p.ID)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

// DeleteProduct removes a product by ID and reports whether a row was deleted.
// Returns (false, nil) if the product did not exist (callers can map to 404).
func (s *Store) DeleteProduct(id int64) (bool, error) {
	const query = `DELETE FROM products WHERE id = ?`
	res, err := s.DB.Exec(query, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}
