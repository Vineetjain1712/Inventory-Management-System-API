# 📦 Inventory Management System API

**Simple, reliable backend API for product & inventory management — written in Go.**
Cleanly layered (`handlers → service → store`), Docker-ready, and designed for straightforward local development or quick deployment.

---

**Repo:** [https://github.com/Vineetjain1712/Inventory-Management-System-API.git](https://github.com/Vineetjain1712/Inventory-Management-System-API.git)

**Docker image:** `vineetjaindev/inventory-api`

---

## 🔍 What this does (TL;DR)

* Full CRUD for products (`id`, `name`, `description`, `stock_quantity`, `low_stock_threshold`)
* Inventory operations: increase / decrease stock (decrease fails on insufficient stock)
* Bonus: list products below `low_stock_threshold`
* SQLite-backed file DB (`inventory.db`) and atomic store updates to avoid simple race conditions
* Makefile for common dev tasks and Docker image for easy runs

---

## ✅ Quick start — run in Docker (recommended)

Run this on any machine with Docker:

```bash
# pull the published image
docker pull vineetjaindev/inventory-api:latest

# create a folder on host to persist the sqlite DB
mkdir -p ~/inventory-data

# run container, mount host folder to /app so inventory.db persists
docker run -d --name inventory-api \
  -p 8080:8080 \
  -v ~/inventory-data:/app \
  vineetjaindev/inventory-api:latest
```

* App will be available at `http://localhost:8080`.
* The app uses `./inventory.db` from its working directory; mounting `~/inventory-data` to `/app` ensures the DB file survives container restarts.

---

## ⚙️ Developer: run locally with Makefile

Clone and use the Makefile targets included in the repository:

```bash
git clone https://github.com/Vineetjain1712/Inventory-Management-System-API.git
cd Inventory-Management-System-API

# download modules
make mod

# run locally (dev)
make run         # runs `go run ./cmd/main.go`

# build a static linux binary
make build       # produces ./inventory-api

# run tests
make test

# docker: build & run locally
make docker-build
make docker-run   # note: make docker-run runs the image without a volume; prefer manual docker run with -v for persistence
```

---

## 🔌 API — endpoints & examples

Base URL: `http://localhost:8080`

### Product endpoints

* `POST /products` — Create product
  Body (JSON):

  ```json
  {
    "name":"Laptop",
    "description":"Gaming laptop",
    "stock_quantity":10,
    "low_stock_threshold":3
  }
  ```
* `GET /products` — List all products
* `GET /products/{id}` — Get a product
* `PUT /products/{id}` — Update a product (send full product JSON)
* `DELETE /products/{id}` — Delete a product

### Inventory endpoints

* `POST /products/{id}/increase` — Increase stock
  Body: `{"quantity": 5}` (or `{"amount": 5}` depending on your client)
* `POST /products/{id}/decrease` — Decrease stock (fails when insufficient)
  Body: `{"quantity": 3}`

### Bonus

* `GET /products/low-stock` — List products where `stock_quantity < low_stock_threshold` (and threshold > 0)

---

## 🔁 Minimal `curl` examples

Create product:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","stock_quantity":10,"low_stock_threshold":3}'
```

List products:

```bash
curl http://localhost:8080/products
```

Increase stock:

```bash
curl -X POST http://localhost:8080/products/1/increase \
  -H "Content-Type: application/json" \
  -d '{"quantity":5}'
```

Decrease stock:

```bash
curl -X POST http://localhost:8080/products/1/decrease \
  -H "Content-Type: application/json" \
  -d '{"quantity":3}'
```

Get product:

```bash
curl http://localhost:8080/products/1
```

Update product:

```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Laptop Pro","description":"High-end","stock_quantity":15,"low_stock_threshold":4}'
```

Delete product:

```bash
curl -X DELETE http://localhost:8080/products/1
```

---

## 📋 Expected status codes (recommended mapping)

* `201 Created` — successful creation
* `200 OK` — successful read/list/update returning JSON
* `204 No Content` — successful update/delete with no body
* `400 Bad Request` — invalid payload, invalid stock operation (e.g., insufficient stock, negative quantities)
* `404 Not Found` — product not found
* `500 Internal Server Error` — unexpected server/DB error

---

## 🧪 Tests

Run unit tests for business logic and store functions:

```bash
make test
# or
go test ./...
```

---

## 🛠️ Notes & recommendations

* The DB file used by the app is `inventory.db` in the working directory. For persistent data when using Docker, mount a host directory to `/app` (as shown above).
* IDs are **assigned by SQLite** (use `INTEGER PRIMARY KEY`). Do not create IDs in application code.
* Atomic DB updates are used for stock adjustments; this reduces race conditions compared to read-modify-write in the app layer.
* If you push a Docker image to Docker Hub, tag it with your username (`vineetjaindev/inventory-api`) and `docker push`.

---

## 🧾 License

This project is available under the **MIT License**. See `LICENSE` in the repository.

---

## 🤝 Contributing

Straightforward process:

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/your-feature`)
3. Add tests for feature or bug fix
4. Open a PR with description and rationale

---

Postman URL : 
https://orange-resonance-433212.postman.co/workspace/My-Workspace~8b91b72a-012d-4a93-88a2-c7240cf47e99/collection/34082492-115b798a-f327-4f52-9ded-2b719246877b?action=share&source=copy-link&creator=34082492
