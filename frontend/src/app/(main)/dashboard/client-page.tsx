"use client";

import { useUser } from "@/context/user-context";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function DashboardClientPage() {
	const { user } = useUser();

	const handleLogout = () => {
		// Create and submit a form to the logout endpoint
		const form = document.createElement("form");
		form.method = "POST";
		form.action = `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/logout`;
		document.body.appendChild(form);
		form.submit();
	};

	return (
		<div className="w-full max-w-4xl mx-auto">
			<Card className="shadow-lg border-0 bg-card/95 backdrop-blur-sm">
				<CardHeader className="space-y-1 pb-6">
					<CardTitle className="text-2xl font-semibold text-center text-card-foreground">
						Dashboard
					</CardTitle>
					<p className="text-sm text-muted-foreground text-center">
						Welcome back!
					</p>
				</CardHeader>
				<CardContent className="space-y-6">
					<div className="text-center">
						<p className="text-lg">
							<span className="font-medium">Email:</span> {user?.email}
						</p>
						{user?.authProvider && (
							<p className="text-sm text-muted-foreground mt-2">
								Signed in via: {user.authProvider}
							</p>
						)}
					</div>
					
					<div className="flex justify-center">
						<Button
							variant="outline"
							onClick={handleLogout}
							className="h-11 text-base font-medium"
						>
							Logout
						</Button>
					</div>
				</CardContent>
			</Card>
		</div>
	);
}
