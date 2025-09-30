export const getGithubAuthorizeUrl = (returnTo?: string) => {
	const searchParams = new URLSearchParams();
	if (returnTo) {
		searchParams.set("return_to", returnTo);
	}

	const apiUrl =
		process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

	return `${apiUrl}/auth/github?${searchParams.toString()}`;
};
