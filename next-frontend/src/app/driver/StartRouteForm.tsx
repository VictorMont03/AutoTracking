"use client";

import { PropsWithChildren, useActionState } from "react";
import { startRouteAction } from "./start-route.action";

export function StartRouteForm(props: PropsWithChildren) {
	const [state, formAction] = useActionState<
		{
			error?: string;
			success?: boolean;
		} | null,
		FormData
	>(startRouteAction, null);
	return (
		<form action={formAction} className="flex flex-col space-y-4">
			{state?.error && (
				<div className="p-4 border rounded text-contrast bg-error">
					{state.error}
				</div>
			)}
			{state?.success && (
				<div className="p-4 border rounded text-contrast bg-success">
					Rota iniciada com sucesso!
				</div>
			)}
			{props.children}
		</form>
	);
}