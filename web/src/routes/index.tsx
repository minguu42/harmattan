import {createFileRoute} from "@tanstack/react-router"
import {PencilIcon, Trash2Icon} from "lucide-react"
import {Menu, MenuItem} from "../components/Menu.tsx"

export const Route = createFileRoute("/")({
	component: RouteComponent,
})

function RouteComponent() {
	return (
		<div className="grid place-items-center">
			<Menu>
				<MenuItem icon={PencilIcon} label="変更" />
				<MenuItem icon={Trash2Icon} label="削除" />
			</Menu>
		</div>
	)
}
