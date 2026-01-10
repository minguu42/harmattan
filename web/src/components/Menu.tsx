import {Menu as BMenu} from "@base-ui/react"
import {EllipsisIcon, type LucideIcon} from "lucide-react"
import type {ReactNode} from "react"
import {IconButton} from "./IconButton.tsx"

type MenuProps = {
	children: ReactNode
}

export function Menu({children}: MenuProps) {
	return (
		<BMenu.Root>
			<BMenu.Trigger render={<IconButton icon={EllipsisIcon} />} />
			<BMenu.Portal>
				<BMenu.Positioner side="bottom" align="end">
					<BMenu.Popup className="w-54 p-1 rounded-2xl bg-surface-container-low shadow-custom">
						{children}
					</BMenu.Popup>
				</BMenu.Positioner>
			</BMenu.Portal>
		</BMenu.Root>
	)
}

type MenuItemProps = {
	icon: LucideIcon
	label: string
	onClick?: () => void
}

export function MenuItem({icon: Icon, label, onClick}: MenuItemProps) {
	return (
		<BMenu.Item className="h-12 flex items-center pl-3 rounded-xl state-layer-parent" onClick={onClick}>
			<div className="state-layer bg-on-surface-variant" />
			<Icon className="size-5 text-on-surface-variant" />
			<div className="w-2" />
			<span className="text-on-surface font-medium">{label}</span>
		</BMenu.Item>
	)
}
