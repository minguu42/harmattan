import {createFileRoute} from "@tanstack/react-router"
import {IconButton} from "./-IconButton.tsx"
import {Pencil} from "lucide-react"

export const Route = createFileRoute("/")({
	component: Index,
})

function Index() {
	return (
		<div className="grid place-items-center min-h-screen">
			<IconButton icon={Pencil} />
		</div>
	)
}
