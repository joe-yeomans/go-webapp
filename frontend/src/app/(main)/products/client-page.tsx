"use client";

import { useEffect, useState } from "react";
import { productApi, Product } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Image from "next/image";

export default function ProductsClientPage() {
	const [products, setProducts] = useState<Product[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const fetchProducts = async () => {
			try {
				setLoading(true);
				const response = await productApi.getProducts();
				setProducts(response.products);
			} catch (err) {
				setError("Failed to load products");
				console.error("Error fetching products:", err);
			} finally {
				setLoading(false);
			}
		};

		fetchProducts();
	}, []);

	const formatPrice = (price: number) => {
		return new Intl.NumberFormat("en-US", {
			style: "currency",
			currency: "USD",
		}).format(price);
	};

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString();
	};

	if (loading) {
		return (
			<div className="w-full max-w-7xl mx-auto">
				<Card>
					<CardHeader>
						<CardTitle>Products</CardTitle>
					</CardHeader>
					<CardContent>
						<div className="flex justify-center items-center h-32">
							<div className="text-muted-foreground">
								Loading products...
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		);
	}

	if (error) {
		return (
			<div className="w-full max-w-7xl mx-auto">
				<Card>
					<CardHeader>
						<CardTitle>Products</CardTitle>
					</CardHeader>
					<CardContent>
						<div className="flex justify-center items-center h-32">
							<div className="text-destructive">{error}</div>
						</div>
					</CardContent>
				</Card>
			</div>
		);
	}

	return (
		<div className="w-full max-w-7xl mx-auto">
			<Card>
				<CardHeader>
					<CardTitle className="text-2xl font-semibold">
						Products
					</CardTitle>
					<p className="text-sm text-muted-foreground">
						Manage and view all products in the system
					</p>
				</CardHeader>
				<CardContent>
					<div className="overflow-x-auto">
						<table className="w-full border-collapse">
							<thead>
								<tr className="border-b">
									<th className="text-left p-3 font-medium">
										Product
									</th>
									<th className="text-left p-3 font-medium">
										Category
									</th>
									<th className="text-left p-3 font-medium">
										Brand
									</th>
									<th className="text-left p-3 font-medium">
										Price
									</th>
									<th className="text-left p-3 font-medium">
										Stock
									</th>
									<th className="text-left p-3 font-medium">
										Rating
									</th>
									<th className="text-left p-3 font-medium">
										Created
									</th>
									<th className="text-left p-3 font-medium">
										Actions
									</th>
								</tr>
							</thead>
							<tbody>
								{products.map((product) => (
									<tr
										key={product.id}
										className="border-b hover:bg-muted/50"
									>
										<td className="p-3">
											<div className="flex items-center space-x-3">
												<div className="flex-shrink-0">
													<Image
														src={product.imageUrl}
														alt={product.name}
														width={48}
														height={48}
														className="h-12 w-12 rounded-lg object-cover"
														unoptimized
													/>
												</div>
												<div className="min-w-0 flex-1">
													<p className="font-medium text-sm truncate">
														{product.name}
													</p>
													<p className="text-xs text-muted-foreground truncate max-w-xs">
														{product.description}
													</p>
												</div>
											</div>
										</td>
										<td className="p-3">
											<span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
												{product.category}
											</span>
										</td>
										<td className="p-3 text-sm text-muted-foreground">
											{product.brand}
										</td>
										<td className="p-3 font-medium">
											{formatPrice(product.price)}
										</td>
										<td className="p-3">
											<div className="flex items-center space-x-2">
												<span
													className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
														product.inStock
															? "bg-green-100 text-green-800"
															: "bg-red-100 text-red-800"
													}`}
												>
													{product.inStock
														? "In Stock"
														: "Out of Stock"}
												</span>
												{product.inStock && (
													<span className="text-xs text-muted-foreground">
														({product.stockCount})
													</span>
												)}
											</div>
										</td>
										<td className="p-3">
											<div className="flex items-center space-x-1">
												<span className="text-sm font-medium">
													{product.rating.toFixed(1)}
												</span>
												<span className="text-xs text-muted-foreground">
													({product.reviewCount})
												</span>
											</div>
										</td>
										<td className="p-3 text-sm text-muted-foreground">
											{formatDate(product.createdAt)}
										</td>
										<td className="p-3">
											<Button
												variant="outline"
												size="sm"
												onClick={() => {
													// TODO: Implement product details modal or page
													console.log(
														"View product:",
														product.id
													);
												}}
											>
												View
											</Button>
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>

					{products.length === 0 && (
						<div className="text-center py-8">
							<p className="text-muted-foreground">
								No products found
							</p>
						</div>
					)}

					<div className="mt-4 text-sm text-muted-foreground">
						Total: {products.length} products
					</div>
				</CardContent>
			</Card>
		</div>
	);
}
