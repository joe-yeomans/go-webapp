package data

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Brand       string    `json:"brand"`
	ImageURL    string    `json:"imageUrl"`
	InStock     bool      `json:"inStock"`
	StockCount  int       `json:"stockCount"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"reviewCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetMockProducts returns a slice of mock products for testing and development
func GetMockProducts() []Product {
	return []Product{
		{
			ID:          "prod_001",
			Name:        "Wireless Bluetooth Headphones",
			Description: "High-quality wireless headphones with noise cancellation and 30-hour battery life",
			Price:       199.99,
			Category:    "Electronics",
			Brand:       "SoundTech",
			ImageURL:    "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500",
			InStock:     true,
			StockCount:  25,
			Rating:      4.5,
			ReviewCount: 128,
			CreatedAt:   time.Now().AddDate(0, -2, -15),
			UpdatedAt:   time.Now().AddDate(0, 0, -5),
		},
		{
			ID:          "prod_002",
			Name:        "Smart Fitness Watch",
			Description: "Advanced fitness tracking with heart rate monitor, GPS, and water resistance up to 50m",
			Price:       299.99,
			Category:    "Electronics",
			Brand:       "FitTech",
			ImageURL:    "https://images.unsplash.com/photo-1544117519-31a4b719223d?w=500",
			InStock:     true,
			StockCount:  18,
			Rating:      4.3,
			ReviewCount: 89,
			CreatedAt:   time.Now().AddDate(0, -1, -20),
			UpdatedAt:   time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          "prod_003",
			Name:        "Organic Cotton T-Shirt",
			Description: "Soft, comfortable organic cotton t-shirt available in multiple colors",
			Price:       29.99,
			Category:    "Clothing",
			Brand:       "EcoWear",
			ImageURL:    "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500",
			InStock:     true,
			StockCount:  45,
			Rating:      4.7,
			ReviewCount: 203,
			CreatedAt:   time.Now().AddDate(0, -3, -10),
			UpdatedAt:   time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          "prod_004",
			Name:        "Ceramic Coffee Mug",
			Description: "Handcrafted ceramic mug perfect for your morning coffee or tea",
			Price:       15.99,
			Category:    "Home & Kitchen",
			Brand:       "Artisan Pottery",
			ImageURL:    "https://images.unsplash.com/photo-1514228742587-6b1558fcf93a?w=500",
			InStock:     true,
			StockCount:  32,
			Rating:      4.8,
			ReviewCount: 156,
			CreatedAt:   time.Now().AddDate(0, -2, -5),
			UpdatedAt:   time.Now().AddDate(0, 0, -3),
		},
		{
			ID:          "prod_005",
			Name:        "Wireless Phone Charger",
			Description: "Fast wireless charging pad compatible with all Qi-enabled devices",
			Price:       49.99,
			Category:    "Electronics",
			Brand:       "PowerUp",
			ImageURL:    "https://images.unsplash.com/photo-1583394838336-acd977736f90?w=500",
			InStock:     false,
			StockCount:  0,
			Rating:      4.2,
			ReviewCount: 67,
			CreatedAt:   time.Now().AddDate(0, -1, -12),
			UpdatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			ID:          "prod_006",
			Name:        "Leather Wallet",
			Description: "Genuine leather wallet with RFID blocking technology and multiple card slots",
			Price:       79.99,
			Category:    "Accessories",
			Brand:       "LeatherCraft",
			ImageURL:    "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500",
			InStock:     true,
			StockCount:  12,
			Rating:      4.6,
			ReviewCount: 94,
			CreatedAt:   time.Now().AddDate(0, -4, -8),
			UpdatedAt:   time.Now().AddDate(0, 0, -4),
		},
		{
			ID:          "prod_007",
			Name:        "Bluetooth Speaker",
			Description: "Portable waterproof Bluetooth speaker with 360-degree sound and 12-hour battery",
			Price:       89.99,
			Category:    "Electronics",
			Brand:       "AudioWave",
			ImageURL:    "https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=500",
			InStock:     true,
			StockCount:  22,
			Rating:      4.4,
			ReviewCount: 112,
			CreatedAt:   time.Now().AddDate(0, -2, -18),
			UpdatedAt:   time.Now().AddDate(0, 0, -6),
		},
		{
			ID:          "prod_008",
			Name:        "Yoga Mat",
			Description: "Non-slip yoga mat made from eco-friendly materials with carrying strap",
			Price:       39.99,
			Category:    "Sports & Fitness",
			Brand:       "ZenFit",
			ImageURL:    "https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=500",
			InStock:     true,
			StockCount:  38,
			Rating:      4.9,
			ReviewCount: 187,
			CreatedAt:   time.Now().AddDate(0, -1, -25),
			UpdatedAt:   time.Now().AddDate(0, 0, -8),
		},
	}
}
