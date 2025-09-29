import type { NextConfig } from "next";

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
};

export default nextConfig;
