export type User = {
	email: string;
	isEmailVerified: boolean;
	googleId?: string;
	githubId?: string;
	name?: string;
	picture?: string;
	authProvider?: string;
};
