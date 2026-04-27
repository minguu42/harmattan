import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"
import {env} from "../env.ts"

type Task = {
	id: string
	name: string
}

type Tasks = {
	tasks: Task[]
}

function isTask(arg: unknown): arg is Task {
	const t = arg as Task
	return typeof t?.id === "string" && typeof t?.name === "string"
}

function isTasks(arg: unknown): arg is Tasks {
	const ts = arg as Tasks
	return Array.isArray(ts?.tasks) && ts?.tasks.every(isTask)
}

export function useTasks(projectID: string) {
	return useQuery({
		queryKey: ["projects", projectID, "tasks"],
		queryFn: async () => {
			const response = await fetch(`${env.apiBaseURL}/projects/${projectID}/tasks?limit=10&offset=0`, {
				method: "GET",
				headers: {"Authorization": `Bearer ${env.apiToken}`},
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
			const data: unknown = await response.json()
			if (isTasks(data)) {
				return data.tasks
			}
			throw new Error("invalid response body")
		},
	})
}

export function useCreateTask(projectID: string) {
	const queryClient = useQueryClient()
	return useMutation({
		mutationFn: async (name: string) => {
			const response = await fetch(`${env.apiBaseURL}/projects/${projectID}/tasks`, {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${env.apiToken}`,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({name}),
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
			return await response.json()
		},
		onSuccess: () => {
			void queryClient.invalidateQueries({queryKey: ["projects", projectID, "tasks"]})
		},
	})
}

export function useCompleteTask(projectID: string) {
	const queryClient = useQueryClient()
	return useMutation({
		mutationFn: async (taskID: string) => {
			const response = await fetch(`${env.apiBaseURL}/tasks/${taskID}`, {
				method: "PATCH",
				headers: {
					"Authorization": `Bearer ${env.apiToken}`,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({completed_at: new Date().toISOString()}),
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
		},
		onSuccess: () => {
			void queryClient.invalidateQueries({queryKey: ["projects", projectID, "tasks"]})
		},
	})
}

export function useDeleteTask(projectID: string) {
	const queryClient = useQueryClient()
	return useMutation({
		mutationFn: async (taskID: string) => {
			const response = await fetch(`${env.apiBaseURL}/tasks/${taskID}`, {
				method: "DELETE",
				headers: {"Authorization": `Bearer ${env.apiToken}`},
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
		},
		onSuccess: () => {
			void queryClient.invalidateQueries({queryKey: ["projects", projectID, "tasks"]})
		},
	})
}
