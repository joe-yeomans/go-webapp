import { getAuthenticatedUser } from "@/utils/user";

export default async function MainLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	const user = await getAuthenticatedUser();

	return (
			<div className="min-h-screen bg-background">
				{/* Header with user info */}
				<header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
					<div className="container flex h-16 items-center justify-between px-4">
						{/* Logo/App name */}
						<div className="flex items-center space-x-2">
							<h1 className="text-xl font-semibold">Your App</h1>
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
									<img
										src={user.picture}
										alt={user.name || user.email}
										className="h-8 w-8 rounded-full object-cover"
									/>
								) : (
									<span className="text-xs font-medium">
										{user?.email?.charAt(0).toUpperCase() || "U"}
									</span>
								)}
							</div>
						</div>
					</div>
				</header>
                                
				{/* Main content */}
				<main className="container px-4 py-6">
					{children}
				</main>
			</div>
	);
}
