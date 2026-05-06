function requireString(name: string): string {
  const value = import.meta.env[name];
  if (typeof value !== "string" || value === "") {
    throw new Error(`environment variable ${name} is required`);
  }
  return value;
}

export const env = {
  apiBaseURL: requireString("VITE_API_BASE_URL"),
  apiToken: requireString("VITE_API_TOKEN"),
} as const;
