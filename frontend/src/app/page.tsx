import SignUp from "@/components/sign-up";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
	return (
		<div className="min-h-screen flex items-center justify-center bg-background">
			<div className="w-full max-w-md">
				<SignUp />
			</div>
		</div>
	)
}
