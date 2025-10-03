# 📦 Inventory Management System API

A simple backend application to manage products and inventory in a warehouse.
This project is designed to demonstrate **clean code, layered architecture, API design, and business logic handling** in Go.

---

## 🚀 Features

* **Product Management**

  * Create, Read, Update, Delete (CRUD) products
  * Product fields: `id`, `name`, `description`, `stock_quantity`, `low_stock_threshold`

* **Inventory Management**

  * Increase stock
  * Decrease stock (with validation: stock cannot go below zero)

* **Bonus**

  * List products that are **below low stock threshold**
  * Dockerized setup for easy run
  * Unit tests for stock logic (edge cases included)

---

## 🗂️ Project Structure

```
inventory-api/
│── cmd/                # Application entrypoint
│    └── main.go
│── internal/
│    ├── models/        # Data models (Product)
│    ├── store/         # In-memory/DB store
│    ├── service/       # Business logic (inventory rules)
│    ├── handlers/      # HTTP handlers (REST APIs)
│── pkg/                # Utilities (error handling, responses)
│── tests/              # Unit tests
│── Dockerfile
│── docker-compose.yml  # (optional for DB)
│── go.mod
│── README.md
```

---

## ⚙️ Setup & Run

### 1. Clone Repository

```bash
git clone https://github.com/<your-username>/inventory-api.git
cd inventory-api
```

### 2. Run Locally (without Docker)

```bash
go run cmd/main.go
```

The server will start at `http://localhost:8080`

### 3. Run with Docker

```bash
docker build -t inventory-api .
docker run -p 8080:8080 inventory-api
```

(Optional: with DB integration, use `docker-compose up`)

---

## 📡 API Endpoints

### Products

* `POST   /products` → Create product
* `GET    /products` → List all products
* `GET    /products/{id}` → Get product by ID
* `PUT    /products/{id}` → Update product
* `DELETE /products/{id}` → Delete product

### Inventory

* `POST /products/{id}/increase` → Increase stock
* `POST /products/{id}/decrease` → Decrease stock (fails if insufficient stock)

### Bonus

* `GET /products/low-stock` → List products below `low_stock_threshold`

---

## 🧪 Running Tests

Run all tests:

```bash
go test ./...
```

---

## 📝 Design Choices & Assumptions

1. **Storage**:

   * Started with an **in-memory store** for simplicity.
   * Can be easily swapped with a persistent DB (e.g., Postgres/SQLite).

2. **Architecture**:

   * Layered: `handlers → service → store`.
   * Keeps business logic separate from HTTP layer.

3. **Error Handling**:

   * Stock operations return `400 Bad Request` if invalid (e.g., insufficient stock).
   * Consistent JSON error responses.

4. **Scalability**:

   * Ready to be extended with authentication, DB, gRPC endpoints, etc.

---

## 🔮 Next Steps (Learning Extensions)

* Add **authentication** (JWT).
* Replace in-memory store with **Postgres + GORM**.
* Add **Swagger docs** for API.
* CI/CD with GitHub Actions.

---

## 👨‍💻 Author

Built as part of a backend learning exercise with **Go, REST APIs, Docker, and Clean Architecture**.
