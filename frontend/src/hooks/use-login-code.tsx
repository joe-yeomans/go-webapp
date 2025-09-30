"use client";

import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { toast } from "sonner";

export default function useLoginCode() {
	const router = useRouter();

	const returnFunc = async (
		email: string,
		returnTo?: string | null
	): Promise<void> => {
		const response = await api.post("/code", { email });
		if (response.status === 200) {
			const params = new URLSearchParams();
			params.set("email", email);
			if (returnTo) {
				params.set("return_to", returnTo);
			}
			router.push(`/verify?${params.toString()}`);
		} else {
			toast.error("Failed to send code");
		}
	};

	return returnFunc;
}
