import {createFileRoute} from "@tanstack/react-router"

export const Route = createFileRoute("/")({
	component: Index,
})

function Index() {
	return (
		<div className="grid place-items-center min-h-screen">
			<h1 className="text-3xl">Hello, World!</h1>
		</div>
	)
}
