"use client";

import { useForm } from "react-hook-form";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useSearchParams } from "next/navigation";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "./ui/form";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import useLoginCode from "@/hooks/use-login-code";
import { FaGoogle } from "react-icons/fa";
import GithubLoginButton from "./github-login-button";

const signUpSchema = z.object({
	email: z.email(),
});

type SignUpSchema = z.infer<typeof signUpSchema>;

export default function SignUp() {
	const [isLoading, setIsLoading] = useState(false);
	const searchParams = useSearchParams();
	const returnTo = searchParams.get("return_to");
	const sendLoginCode = useLoginCode();
	const form = useForm<SignUpSchema>({
		resolver: zodResolver(signUpSchema),
		defaultValues: {
			email: "",
		},
	});

	const onSubmit = async (data: SignUpSchema) => {
		setIsLoading(true);
		await sendLoginCode(data.email, returnTo);
		setIsLoading(false);
	};

	const handleGoogleLogin = () => {
		// Redirect to backend Google OAuth endpoint with return_to parameter
		const apiUrl =
			process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";
		const googleUrl = `${apiUrl}/auth/google`;

		if (returnTo) {
			window.location.href = `${googleUrl}?return_to=${encodeURIComponent(
				returnTo
			)}`;
		} else {
			window.location.href = googleUrl;
		}
	};

	return (
		<div className="w-full max-w-md mx-auto">
			<Card className="shadow-lg border-0 backdrop-blur-sm">
				<CardHeader className="space-y-1 pb-6">
					<CardTitle className="text-2xl font-semibold text-center">
						Sign Up
					</CardTitle>
					<p className="text-sm text-muted-foreground text-center">
						Enter your email to get started
					</p>
				</CardHeader>
				<CardContent className="space-y-6">
					{/* OAuth buttons - Primary options */}
					<div className="space-y-3">
						<Button
							type="button"
							onClick={handleGoogleLogin}
							className="w-full h-11 text-base font-medium"
						>
							<FaGoogle className="mr-2 h-4 w-4" />
							Continue with Google
						</Button>

						<GithubLoginButton returnTo={returnTo ?? undefined} />
					</div>

					<div className="relative">
						<div className="absolute inset-0 flex items-center">
							<span className="w-full border-t" />
						</div>
						<div className="relative flex justify-center text-xs uppercase">
							<span className="bg-background px-2 text-muted-foreground">
								Or sign in with email
							</span>
						</div>
					</div>

					{/* Email form - Secondary option */}
					<Form {...form}>
						<form
							onSubmit={form.handleSubmit(onSubmit)}
							className="space-y-6"
						>
							<FormField
								control={form.control}
								name="email"
								render={({ field }) => (
									<FormItem className="space-y-2">
										<FormLabel className="text-sm font-medium text-card-foreground">
											Email Address
										</FormLabel>
										<FormControl>
											<Input
												{...field}
												type="email"
												placeholder="Enter your email"
												className="h-11 px-4 text-base"
											/>
										</FormControl>
										<FormMessage className="text-sm" />
									</FormItem>
								)}
							/>
							<Button
								type="submit"
								disabled={isLoading}
								variant="secondary"
								className="w-full h-11 text-base font-medium disabled:opacity-50 disabled:cursor-not-allowed"
							>
								{isLoading
									? "Sending Code..."
									: "Send Verification Code"}
							</Button>
						</form>
					</Form>

					<div className="text-center">
						<p className="text-xs text-muted-foreground">
							By signing up, you agree to our Terms of Service and
							Privacy Policy
						</p>
					</div>
				</CardContent>
			</Card>
		</div>
	);
}
