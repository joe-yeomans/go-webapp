"use client";

import { useForm } from "react-hook-form";
import { useSearchParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "./ui/form";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "./ui/input-otp";
import { Button } from "./ui/button";

const verifyCodeSchema = z.object({
	code: z
		.string()
		.min(6, "Code must be at least 6 characters")
		.max(6, "Code must be exactly 6 characters"),
});

type VerifyCodeSchema = z.infer<typeof verifyCodeSchema>;

export default function VerifyCode() {
	const searchParams = useSearchParams();
	const email = searchParams.get("email") || "";
	const error = searchParams.get("error") || "";

	// Use environment variable or default to localhost:8080
	const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

	const form = useForm<VerifyCodeSchema>({
		resolver: zodResolver(verifyCodeSchema),
		defaultValues: {
			code: "",
		},
	});

	return (
		<Card className="shadow-lg border-0 bg-card/95 backdrop-blur-sm">
				<CardHeader className="space-y-1 pb-6">
					<CardTitle className="text-2xl font-semibold text-center text-card-foreground">
						Verify Code
					</CardTitle>
					<p className="text-sm text-muted-foreground text-center">
						Enter the 6-digit code sent to your email
					</p>
					{error && (
						<div className="bg-destructive/10 border border-destructive/20 rounded-md p-3 mt-4">
							<p className="text-sm text-destructive text-center">
								{decodeURIComponent(error)}
							</p>
						</div>
					)}
				</CardHeader>
				<CardContent className="space-y-6">
					<Form {...form}>
						<form
							action={`${apiUrl}/api/verify`}
							method="POST"
							className="space-y-6"
						>
							{/* Hidden email field */}
							<input type="hidden" name="email" value={email} />

							<FormField
								control={form.control}
								name="code"
								render={({ field }) => (
									<FormItem className="space-y-4">
										<FormLabel className="text-sm font-medium text-card-foreground text-center block">
											Verification Code
										</FormLabel>
										<FormControl>
											<div className="flex justify-center">
												<InputOTP maxLength={6} {...field}>
													<InputOTPGroup className="gap-2">
														<InputOTPSlot index={0} className="h-12 w-12 text-lg" />
														<InputOTPSlot index={1} className="h-12 w-12 text-lg" />
														<InputOTPSlot index={2} className="h-12 w-12 text-lg" />
														<InputOTPSlot index={3} className="h-12 w-12 text-lg" />
														<InputOTPSlot index={4} className="h-12 w-12 text-lg" />
														<InputOTPSlot index={5} className="h-12 w-12 text-lg" />
													</InputOTPGroup>
												</InputOTP>
											</div>
										</FormControl>
										<FormMessage className="text-sm text-center" />
									</FormItem>
								)}
							/>
							<Button
								type="submit"
								className="w-full h-11 text-base font-medium"
							>
								Verify Code
							</Button>
						</form>
					</Form>
					<div className="text-center">
						<p className="text-xs text-muted-foreground">
							Didn&apos;t receive a code?{" "}
							<button
								type="button"
								className="text-primary hover:underline font-medium"
								onClick={() =>
									console.log("Resend code requested")
								}
							>
								Resend
							</button>
						</p>
					</div>
				</CardContent>
		</Card>
	);
}
