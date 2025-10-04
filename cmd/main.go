package main

import (
	"log"
	"net/http"

	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/handlers"
	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/service"
	"github.com/Vineetjain1712/Inventory-Management-System-API/internal/store"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize store (SQLite)
	db, err := store.NewStore("inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize service
	svc := service.NewService(db)

	// Initialize handlers
	h := handlers.NewProductHandler(svc)

	// Setup router
	r := mux.NewRouter()

	// Bonus: low stock
	r.HandleFunc("/products/low-stock", h.ListLowStockProducts).Methods("GET")

	// CRUD endpoints
	r.HandleFunc("/products", h.CreateProduct).Methods("POST")
	r.HandleFunc("/products", h.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", h.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")

	// Stock endpoints
	r.HandleFunc("/products/{id}/increase", h.IncreaseStock).Methods("POST")
	r.HandleFunc("/products/{id}/decrease", h.DecreaseStock).Methods("POST")

	// Start server
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
