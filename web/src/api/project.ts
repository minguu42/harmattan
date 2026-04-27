import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"
import {env} from "../env.ts"

export type Project = {
	id: string
	name: string
}

type Projects = {
	projects: Project[]
}

function isProject(arg: unknown): arg is Project {
	const p = arg as Project
	return typeof p?.id === "string" && typeof p?.name === "string"
}

function isProjects(arg: unknown): arg is Projects {
	const ps = arg as Projects
	return Array.isArray(ps?.projects) && ps?.projects.every(isProject)
}

export function useProject(projectID: string) {
	return useQuery({
		queryKey: ["projects", projectID],
		queryFn: async () => {
			const response = await fetch(`${env.apiBaseURL}/projects/${projectID}`, {
				method: "GET",
				headers: {"Authorization": `Bearer ${env.apiToken}`},
			});
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`);
			}
			const data: unknown = await response.json();
			if (isProject(data)) {
				return data;
			}
			throw new Error("invalid response body");
		},
	});
}

export function useProjects() {
	return useQuery({
		queryKey: ["projects"],
		queryFn: async () => {
			const response = await fetch(`${env.apiBaseURL}/projects?limit=10&offset=0`, {
				method: "GET",
				headers: {"Authorization": `Bearer ${env.apiToken}`},
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
			const data: unknown = await response.json()
			if (isProjects(data)) {
				return data.projects
			}
			throw new Error("invalid response body")
		},
	})
}

export function useCreateProject() {
	const c = useQueryClient()
	return useMutation({
		mutationFn: async (name: string) => {
			const response = await fetch(`${env.apiBaseURL}/projects`, {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${env.apiToken}`,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({name, color: "default"}),
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
			return await response.json()
		},
		onSuccess: () => {
			void c.invalidateQueries({queryKey: ["projects"]})
		},
	})
}

export function useDeleteProject() {
	const client = useQueryClient()
	return useMutation({
		mutationFn: async (projectID: string) => {
			const response = await fetch(`${env.apiBaseURL}/projects/${projectID}`, {
				method: "DELETE",
				headers: {"Authorization": `Bearer ${env.apiToken}`},
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
		},
		onSuccess: () => {
			void client.invalidateQueries({queryKey: ["projects"]})
		}
	})
}
