import type { NextConfig } from "next";

// Build CSP as a single-line header value to avoid invalid characters (no newlines/tabs)
const csp = [
	"default-src 'self'",
	"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
	"style-src 'self' 'unsafe-inline'",
	"img-src 'self' https://images.unsplash.com https://avatars.githubusercontent.com",
	"font-src 'self'",
	"connect-src 'self' http://localhost:8080",
	"frame-src 'self'",
	"frame-ancestors 'none'",
	"base-uri 'self'",
	"form-action 'self' http://localhost:8080",
].join("; ");

const nextConfig: NextConfig = {
	async redirects() {
		return [
			{
				source: "/",
				destination: "/dashboard",
				has: [
					{
						type: "cookie",
						key: "session",
					},
				],
				permanent: true,
			},
		];
	},

	async headers() {
		return [
			{
				source: "/(.*)",
				headers: [
					{ key: "Content-Security-Policy", value: csp },
					{ key: "X-Content-Type-Options", value: "nosniff" },
					{
						key: "Referrer-Policy",
						value: "strict-origin-when-cross-origin",
					},
				],
			},
			// far future cache for nextjs build assets
			{
				source: "/_next/static/:path*",
				headers: [
					{
						key: "Cache-Control",
						value: "max-age=31536000, immutable",
					},
				],
			},
		];
	},
};

export default nextConfig;
