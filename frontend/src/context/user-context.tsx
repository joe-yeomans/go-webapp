"use client";

import { User } from "@/schemas/user";
import { createContext, useContext, useState, ReactNode } from "react";

interface UserContextType {
	user: User | null;
	setUser: (user: User | null) => void;
}

const UserContext = createContext<UserContextType>({
	user: null,
	setUser: () => {},
});

export function UserProvider({
	user: _user,
	children,
}: {
	user: User | null;
	children: ReactNode;
}) {
	const [user, setUser] = useState<User | null>(_user);

	return (
		<UserContext.Provider value={{ user, setUser }}>
			{children}
		</UserContext.Provider>
	);
}

export function useUser() {
	return useContext(UserContext);
}
