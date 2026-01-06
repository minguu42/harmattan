import {Dialog, Field, Form} from "@base-ui/react"
import {useState} from "react"
import {IconButton} from "./-IconButton.tsx"
import {
	FolderOpenDotIcon,
	FolderPlusIcon,
	ListTodoIcon,
	type LucideIcon,
	MenuIcon,
	SunIcon,
	Trash2Icon,
} from "lucide-react"
import {Link} from "@tanstack/react-router"
import {type Project, useCreateProject, useDeleteProject, useProjects} from "./-api_project.ts"
import {Button} from "./-Button.tsx"

export function NavigationDrawer() {
	const [open, setOpen] = useState(false)

	function toggleDrawer() {
		setOpen((prev) => !prev)
	}

	return (
		<Dialog.Root open={open} onOpenChange={setOpen}>
			<IconButton icon={MenuIcon} onClick={toggleDrawer} />
			<Dialog.Portal>
				<Dialog.Backdrop className="fixed inset-0 bg-scrim" />
				<Dialog.Popup className="group/nd fixed top-0 left-0 w-90 min-h-screen bg-surface rounded-r-xl px-3">
					<div className="fixed inset-0 bg-scrim hidden group-data-nested-dialog-open/nd:block" />
					<div className="h-16 flex items-center pl-4">
						<p className="text-on-surface-variant">タスク</p>
					</div>
					<ul>
						<Indicator icon={SunIcon} label="今日のタスク" />
						<Indicator icon={ListTodoIcon} label="すべてのタスク" />
						<Indicator icon={Trash2Icon} label="ごみ箱" />
					</ul>
					<div className="border w-78 mx-auto my-px border-on-surface-variant" />
					<div className="h-16 flex items-center pl-4">
						<p className="text-on-surface-variant">プロジェクト</p>
						<div className="flex-1" />
						<ProjectCreateDialog />
					</div>
					<ProjectIndicatorList />
				</Dialog.Popup>
			</Dialog.Portal>
		</Dialog.Root>
	)
}

type IndicatorProps = {
	icon: LucideIcon
	label: string
}

function Indicator({icon: Icon, label}: IndicatorProps) {
	return (
		<li>
			<Link to="/" className="state-layer-parent flex items-center h-14 pl-4 rounded-xl">
				<div className="state-layer bg-on-surface-variant" />
				<div className="state-layer-ring" />
				<Icon className="text-on-surface-variant" />
				<div className="w-3" />
				<p className="text-sm text-on-surface-variant">{label}</p>
			</Link>
		</li>
	)
}

function ProjectCreateDialog() {
	const [open, setOpen] = useState(false)

	return (
		<Dialog.Root open={open} onOpenChange={setOpen}>
			<IconButton icon={FolderPlusIcon} onClick={() => setOpen(true)} />
			<Dialog.Portal>
				<ProjectCreatePopup closeDialog={() => setOpen(false)} />
			</Dialog.Portal>
		</Dialog.Root>
	)
}

type ProjectCreatePopupProps = {
	closeDialog: () => void
}

function ProjectCreatePopup({closeDialog}: ProjectCreatePopupProps) {
	const [name, setName] = useState("")
	const createProject = useCreateProject()

	function addProject() {
		if (name.trim() === "") return
		createProject.mutate(name.trim(), {
			onSuccess: () => setName(""),
		})
		closeDialog()
	}

	return (
		<Dialog.Popup className="fixed top-1/2 left-1/2 w-100 bg-surface -translate-x-1/2 -translate-y-1/2 rounded-xl -mt-8 p-6">
			<Dialog.Title className="text-2xl text-on-surface">プロジェクト作成</Dialog.Title>
			<div className="h-6" />
			<Form onSubmit={(e) => {
				e.preventDefault()
				addProject()
			}}>
				<Field.Root>
					<Field.Control
						className=" w-88 h-14 px-4 text-on-surface placeholder:text-on-surface-variant"
						required
						placeholder="プロジェクト名"
						value={name}
						onChange={(e) => setName(e.target.value)}
					/>
				</Field.Root>
				<div className="flex">
					<div className="flex-1" />
					<Button type="submit" label="作成" />
				</div>
			</Form>
		</Dialog.Popup>
	)
}

function ProjectIndicatorList() {
	const {data: projects, error, isPending, isError} = useProjects()

	if (isPending) {
		return <span>Loading...</span>
	}
	if (isError) {
		return <span>Error: {error.message}</span>
	}
	return (
		<ul>
			{projects.map((p) => <ProjectIndicator key={p.id} project={p} />)}
		</ul>
	)
}

type ProjectIndicatorProps = {
	project: Project
}

function ProjectIndicator({project}: ProjectIndicatorProps) {
	const deleteProject = useDeleteProject()

	return (
		<li>
			<Link
				to="/projects/$projectID/tasks"
				params={{projectID: project.id}}
				className="state-layer-parent group/i flex items-center h-14 pl-4 rounded-xl"
			>
				<div className="state-layer bg-on-surface-variant" />
				<div className="state-layer-ring" />
				<FolderOpenDotIcon className="text-on-surface-variant" />
				<div className="w-3" />
				<div className="test-sm text-on-surface-variant">{project.name}</div>
				<div className="flex-1" />
				<div className="hidden group-hover/i:block">
					<IconButton icon={Trash2Icon} onClick={() => deleteProject.mutate(project.id)} />
				</div>
			</Link>
		</li>
	)
}
