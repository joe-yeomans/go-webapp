package handler

import (
	"encoding/json"
	"go-webapp/internal/data"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Product API handlers

// ProductListResponse represents the response for product list endpoints
type ProductListResponse struct {
	Products []data.Product `json:"products"`
	Total    int            `json:"total"`
}

// handleGetProducts returns all products
func (h *Handler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Handle GET /api/products endpoint running")

	products := data.GetMockProducts()

	response := ProductListResponse{
		Products: products,
		Total:    len(products),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleGetProduct returns a single product by ID
func (h *Handler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Handle GET /api/products/{id} endpoint running")

	productID := chi.URLParam(r, "id")

	products := data.GetMockProducts()

	// Find the product with the given ID
	for _, product := range products {
		if product.ID == productID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	// Product not found
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Product not found",
	})
}

// handleGetProductsByCategory returns products filtered by category
func (h *Handler) HandleGetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Handle GET /api/products/category/{category} endpoint running")

	category := chi.URLParam(r, "category")

	products := data.GetMockProducts()
	var filteredProducts []data.Product

	// Filter products by category (case-insensitive)
	for _, product := range products {
		if strings.EqualFold(product.Category, category) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	response := ProductListResponse{
		Products: filteredProducts,
		Total:    len(filteredProducts),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
