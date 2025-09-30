import axios, { AxiosInstance } from "axios";

// Create axios instance with base configuration
const api: AxiosInstance = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
	timeout: 10000,
	headers: {
		"Content-Type": "application/json",
	},
	withCredentials: true,
});

// Product types
export interface Product {
	id: string;
	name: string;
	description: string;
	price: number;
	category: string;
	brand: string;
	imageUrl: string;
	inStock: boolean;
	stockCount: number;
	rating: number;
	reviewCount: number;
	createdAt: string;
	updatedAt: string;
}

export interface ProductListResponse {
	products: Product[];
	total: number;
}

// Product API functions
export const productApi = {
	// Get all products
	getProducts: async (): Promise<ProductListResponse> => {
		const response = await api.get<ProductListResponse>("/products");
		return response.data;
	},

	// Get a single product by ID
	getProduct: async (id: string): Promise<Product> => {
		const response = await api.get<Product>(`/products/${id}`);
		return response.data;
	},

	// Get products by category
	getProductsByCategory: async (
		category: string
	): Promise<ProductListResponse> => {
		const response = await api.get<ProductListResponse>(
			`/products/category/${encodeURIComponent(category)}`
		);
		return response.data;
	},
};

// Export the axios instance
export default api;
