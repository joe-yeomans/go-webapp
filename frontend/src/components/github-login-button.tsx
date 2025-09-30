import Link from "next/link";
import { Button } from "./ui/button";
import { FaGithub } from "react-icons/fa";
import { getGithubAuthorizeUrl } from "@/utils/auth";

interface Props {
	returnTo?: string;
}

export default function GithubLoginButton({ returnTo }: Props) {
	const url = getGithubAuthorizeUrl(returnTo);

	return (
		<Link href={url}>
			<Button
				variant="outline"
				type="button"
				className="w-full h-11 text-base font-medium"
			>
				<FaGithub className="mr-2 h-4 w-4" />
				Continue with GitHub
			</Button>
		</Link>
	);
}
