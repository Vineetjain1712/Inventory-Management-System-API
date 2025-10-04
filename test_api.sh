# How to Run : 
# chmod +x test_script.sh
# ./test_script.sh



#!/usr/bin/env bash

# A very simple dev test script for your Inventory API.
# It will:
# 1) Create two products
# 2) List all products
# 3) Increase & decrease stock
# 4) Try an invalid decrease (should fail)
# 5) Get a product by ID
# 6) List low-stock products
# 7) Update a product
# 8) Delete a product
# 9) Show final list
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=============================="
echo "1️⃣ Create product: Laptop"
echo "=============================="
curl -X POST "$BASE_URL/products" \
-H "Content-Type: application/json" \
-d '{
  "name": "Laptop",
  "description": "Gaming laptop",
  "stock_quantity": 10,
  "low_stock_threshold": 3
}'
echo -e "\n"

echo "=============================="
echo "2️⃣ Create product: Keyboard"
echo "=============================="
curl -X POST "$BASE_URL/products" \
-H "Content-Type: application/json" \
-d '{
  "name": "Keyboard",
  "description": "Mechanical keyboard",
  "stock_quantity": 5,
  "low_stock_threshold": 2
}'
echo -e "\n"

echo "=============================="
echo "3️⃣ List all products"
echo "=============================="
curl "$BASE_URL/products"
echo -e "\n"

echo "=============================="
echo "4️⃣ Increase stock of product ID 1 by 5"
echo "=============================="
curl -X POST "$BASE_URL/products/1/increase" \
-H "Content-Type: application/json" \
-d '{"quantity": 5}'
echo -e "\n"

echo "=============================="
echo "5️⃣ Decrease stock of product ID 2 by 3"
echo "=============================="
curl -X POST "$BASE_URL/products/2/decrease" \
-H "Content-Type: application/json" \
-d '{"quantity": 3}'
echo -e "\n"

echo "=============================="
echo "6️⃣ Attempt to decrease stock of product ID 2 by 5 (should fail)"
echo "=============================="
curl -X POST "$BASE_URL/products/2/decrease" \
-H "Content-Type: application/json" \
-d '{"quantity": 5}'
echo -e "\n"

echo "=============================="
echo "7️⃣ Get product by ID 1"
echo "=============================="
curl "$BASE_URL/products/1"
echo -e "\n"

echo "=============================="
echo "8️⃣ Get low-stock products"
echo "=============================="
curl "$BASE_URL/products/low-stock"
echo -e "\n"

echo "=============================="
echo "9️⃣ Update product ID 1"
echo "=============================="
curl -X PUT "$BASE_URL/products/1" \
-H "Content-Type: application/json" \
-d '{
  "id": 1,
  "name": "Laptop Pro",
  "description": "High-end gaming laptop",
  "stock_quantity": 15,
  "low_stock_threshold": 4
}'
echo -e "\n"

echo "=============================="
echo "🔟 Delete product ID 2"
echo "=============================="
curl -X DELETE "$BASE_URL/products/2"
echo -e "\n"

echo "=============================="
echo "✅ Final product list"
echo "=============================="
curl "$BASE_URL/products"
echo -e "\n"

echo "=============================="
echo "✅ All tests completed!"
echo "=============================="
