import axios, { AxiosInstance } from "axios";

// Create axios instance with base configuration
const api: AxiosInstance = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
	timeout: 10000,
	headers: {
		"Content-Type": "application/json",
	},
});

// Export the axios instance
export default api;
