import {createFileRoute} from "@tanstack/react-router"
import {IconButton} from "./-IconButton.tsx"
import {Pencil} from "lucide-react"

export const Route = createFileRoute("/")({
	component: RouteComponent,
})

function RouteComponent() {
	return (
		<div className="grid place-items-center">
			<IconButton icon={Pencil} />
		</div>
	)
}
