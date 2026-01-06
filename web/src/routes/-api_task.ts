import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"

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

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMUtBM01QSk1URzRTNVlTTTNXME5TRlk3OSIsImV4cCI6MTc3MDk4MzE0MCwiaWF0IjoxNzYzMjA3MTQwfQ.tJ8WOl0vp3ccLTXdO6bzW5V7CAIkfkw5WU1mKNihIQY"

export function useTasks(projectID: string) {
	return useQuery({
		queryKey: ["projects", projectID, "tasks"],
		queryFn: async () => {
			const response = await fetch(`http://127.0.0.1:8080/projects/${projectID}/tasks?limit=10&offset=0`, {
				method: "GET",
				headers: {"Authorization": `Bearer ${token}`},
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
			const response = await fetch(`http://127.0.0.1:8080/projects/${projectID}/tasks`, {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({name}),
			})
			if (!response.ok) {
				throw new Error(`HTTP error status: ${response.status}`)
			}
			return await response.json() as Task
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
			const response = await fetch(`http://127.0.0.1:8080/tasks/${taskID}`, {
				method: "PATCH",
				headers: {
					"Authorization": `Bearer ${token}`,
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
			const response = await fetch(`http://127.0.0.1:8080/tasks/${taskID}`, {
				method: "DELETE",
				headers: {"Authorization": `Bearer ${token}`},
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
