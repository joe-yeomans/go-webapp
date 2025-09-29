import VerifyCode from "@/components/verify-code";

export default function VerifyPage() {
	return (
		<div className="min-h-screen flex items-center justify-center bg-background">
			<div className="w-full max-w-md">
				<VerifyCode />
			</div>
		</div>
	);
}
