"use client";

import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { toast } from "sonner";

export default function useLoginCode() {
	const router = useRouter();

	const returnFunc = async (email: string): Promise<void> => {
		const response = await api.post("/code", { email });
		if (response.status === 200) {
			router.push(`/verify?email=${encodeURIComponent(email)}`);
		} else {
			toast.error("Failed to send code");
		}
	};

	return returnFunc;
}
