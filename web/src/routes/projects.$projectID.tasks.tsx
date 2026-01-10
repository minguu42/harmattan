import {createFileRoute} from "@tanstack/react-router"
import {useProject} from "../api/project.ts"
import {CheckIcon, CircleIcon, FolderOpenDotIcon, PlusIcon, Trash2Icon} from "lucide-react"
import {useState} from "react"
import {useCompleteTask, useCreateTask, useDeleteTask, useTasks} from "../api/task.ts"
import {Field, Form} from "@base-ui/react"
import {IconButton} from "../components/IconButton.tsx"

export const Route = createFileRoute("/projects/$projectID/tasks")({
	component: RouteComponent,
})

function RouteComponent() {
	const {projectID} = Route.useParams()
	return (
		<div className="max-w-178 mx-auto">
			<div className="h-9" />
			<TaskListHeader projectID={projectID} />
			<div className="h-4" />
			<TaskAddField projectID={projectID} />
			<TaskList projectID={projectID} />
		</div>
	)
}

function TaskListHeader({projectID}: { projectID: string }) {
	const {data: project, error, isPending, isError} = useProject(projectID)

	if (isPending) {
		return <span>Loading...</span>
	}
	if (isError) {
		return <span>Error: {error.message}</span>
	}
	return (
		<div className="flex items-center">
			<FolderOpenDotIcon size={36} className="text-on-surface" />
			<h1 className="ml-1 text-[32px] text-on-surface">{project.name}</h1>
		</div>
	)
}

function TaskAddField({projectID}: { projectID: string }) {
	const [name, setName] = useState("");
	const createTask = useCreateTask(projectID);

	function addTask() {
		if (name.trim() === "") return;
		createTask.mutate(name.trim(), {
			onSuccess: () => setName(""),
		});
	}

	return (
		<div className="p-1 h-10 flex items-center">
			<PlusIcon className="text-[#8e8a90]" />
			<Form
				className="ml-2 w-full"
				onSubmit={(e) => {
					e.preventDefault();
					addTask();
				}}
			>
				<Field.Root>
					<Field.Control
						className="w-full p-1 text-on-surface placeholder:text-[#8e8a90]"
						required
						placeholder="タスク名"
						value={name}
						onChange={(e) => setName(e.target.value)}
					/>
				</Field.Root>
			</Form>
		</div>
	);
}

function TaskList({projectID}: { projectID: string }) {
	const {data: tasks, error, isPending, isError} = useTasks(projectID);
	const completeTask = useCompleteTask(projectID);
	const deleteTask = useDeleteTask(projectID);

	if (isPending) {
		return <span>Loading...</span>;
	}
	if (isError) {
		return <span>Error: {error.message}</span>;
	}
	return (
		<ul className="flex flex-col gap-px">
			{tasks.map((t) => (
				<li key={t.id} className="flex items-center gap-3">
					<IconButton icon={CircleIcon} hoverIcon={CheckIcon} onClick={() => completeTask.mutate(t.id)} />
					<p>{t.name}</p>
					<div className="flex-1" />
					<IconButton icon={Trash2Icon} onClick={() => deleteTask.mutate(t.id)} />
				</li>
			))}
		</ul>
	);
}
