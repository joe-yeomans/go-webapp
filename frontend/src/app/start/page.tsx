import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function StartPage() {
	// Add a small delay to ensure cookie is processed
	await new Promise(resolve => setTimeout(resolve, 1000));

    const requestCookies = await cookies();
    console.log(requestCookies);
	
	// Server-side redirect to dashboard
	redirect("/dashboard");
}
