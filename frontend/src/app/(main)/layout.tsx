import { getAuthenticatedUser } from "@/utils/user";
import Image from "next/image";
import Link from "next/link";

export default async function MainLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	const user = await getAuthenticatedUser();

	return (
		<div className="min-h-screen w-full bg-background">
			{/* Header with user info */}
			<header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
				<div className="flex h-16 w-full items-center justify-between px-4">
					{/* Logo/App name and navigation */}
					<div className="flex items-center space-x-6">
						<h1 className="text-xl font-semibold">Your App</h1>
						<nav className="hidden md:flex items-center space-x-4">
							<Link
								href="/dashboard"
								className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
							>
								Dashboard
							</Link>
							<Link
								href="/products"
								className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
							>
								Products
							</Link>
						</nav>
					</div>

					{/* User info in top right */}
					<div className="flex items-center space-x-3">
						{user?.name && (
							<span className="text-sm text-muted-foreground">
								{user.name}
							</span>
						)}

						{/* Avatar */}
						<div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground">
							{user?.picture ? (
								<Image
									src={user.picture}
									alt={user.name || user.email}
									width={32}
									height={32}
									className="h-8 w-8 rounded-full object-cover"
									unoptimized
									priority
								/>
							) : (
								<span className="text-xs font-medium">
									{user?.email?.charAt(0).toUpperCase() ||
										"U"}
								</span>
							)}
						</div>
					</div>
				</div>
			</header>

			{/* Main content */}
			<main className="w-full px-4 py-6">{children}</main>
		</div>
	);
}
