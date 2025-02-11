"use server";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function startRouteAction(state: any, formData: FormData) {
	const { route_id } = Object.fromEntries(formData);

	if (!route_id) {
		return { error: "Route ID is required" };
	}

	console.log(`Calling start route with route_id: ${route_id}`);


	const response = await fetch(
		`${process.env.NEST_API_URL}/routes/${route_id}/start`,
		{
			method: "POST",
		}
	);

	if (!response.ok) {
		console.error(await response.text());
		return { error: "Failed to start route" };
	}

	return { success: true };
}