import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"

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

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMUtBM01QSk1URzRTNVlTTTNXME5TRlk3OSIsImV4cCI6MTc3MDk4MzE0MCwiaWF0IjoxNzYzMjA3MTQwfQ.tJ8WOl0vp3ccLTXdO6bzW5V7CAIkfkw5WU1mKNihIQY"

export function useProjects() {
	return useQuery({
		queryKey: ["projects"],
		queryFn: async () => {
			const response = await fetch("http://127.0.0.1:8080/projects?limit=10&offset=0", {
				method: "GET",
				headers: {"Authorization": `Bearer ${token}`},
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
			const response = await fetch("http://127.0.0.1:8080/projects", {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
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
