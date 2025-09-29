import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { User } from "./schemas/user";
import api from "./lib/api";

const protectedRoutes = ["/dashboard"];

const isProtectedRoute = (pathname: string) =>
	protectedRoutes.some((route) => pathname.startsWith(route));

export async function middleware(request: NextRequest) {
	const { pathname } = request.nextUrl;

	const hasSessionCookie = !!request.cookies.get("session")?.value;
	console.log(request.cookies.get("session"))

	let user: User | null = null;
	if (hasSessionCookie) {
		try {
			const response = await api.get<User>("/me", {
				headers: {
					Cookie: request.cookies.toString(),
				},
			});
			user = response.data;
		} catch {
			user = null;
		}
	}

	if (isProtectedRoute(pathname) && !user) {
		const url = request.nextUrl.clone();
		url.pathname = "/login";
		return NextResponse.redirect(url);
	}

	const headers = user ? { "x-user": JSON.stringify(user) } : undefined;
	return NextResponse.next({ headers });
}

export const config = {
	matcher: [
		/*
		 * Match all request paths except for the ones starting with:
		 * - api (API routes)
		 * - _next/static (static files)
		 * - _next/image (image optimization files)
		 * - favicon.ico (favicon file)
		 * - public files
		 */
		"/((?!api|_next/static|_next/image|favicon.ico|.*\\.).*)",
	],
};
