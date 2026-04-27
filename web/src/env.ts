function requireString(name: string, value: unknown): string {
	if (typeof value !== "string" || value === "") {
		throw new Error(`environment variable ${name} is required`)
	}
	return value
}

export const env = {
	apiBaseURL: requireString("VITE_API_BASE_URL", import.meta.env.VITE_API_BASE_URL),
	apiToken: requireString("VITE_API_TOKEN", import.meta.env.VITE_API_TOKEN),
} as const
