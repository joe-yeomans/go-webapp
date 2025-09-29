import { User } from "@/schemas/user";
import { headers } from "next/headers";

export const getAuthenticatedUser = async (): Promise<User | null> => {
	const userHeaders = await headers();
	const user = userHeaders.get("x-user");
	if (!user) return null;
	return JSON.parse(user);
};
