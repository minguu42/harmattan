import {createFileRoute} from "@tanstack/react-router"

export const Route = createFileRoute("/projects/$projectID/tasks")({
	component: RouteComponent,
})

function RouteComponent() {
	const {projectID} = Route.useParams()
	return <h1>Hello, プロジェクト{projectID}!</h1>
}
