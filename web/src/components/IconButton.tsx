import type {LucideIcon} from "lucide-react"
import {Button} from "@base-ui/react"

type Props = {
	icon: LucideIcon
	hoverIcon?: LucideIcon
	size?: "xs" | "sm"
	onClick?: () => void
}

const containerSizes = {
	"xs": "size-8",
	"sm": "size-10",
}

const iconSizes = {
	"xs": "size-5",
	"sm": "size-6",
}

export function IconButton({icon: Icon, hoverIcon: HoverIcon, size = "sm", onClick}: Props) {
	return (
		<Button className="group/ib relative size-12 grid place-items-center rounded-xl focus-visible:outline-none" onClick={onClick}>
			<div className={`absolute inset-0 m-auto rounded-[inherit] opacity-0 bg-on-surface-variant group-hover/ib:opacity-10 group-focus-visible/ib:opacity-10 group-active/ib:opacity-12 ${containerSizes[size]}`} />
			<div className={`absolute inset-0 m-auto rounded-[inherit] group-focus-visible/ib:outline-2 group-focus-visible/ib:outline-focus-ring ${containerSizes[size]}`} />
			<Icon className={`${iconSizes[size]} text-on-surface-variant ${HoverIcon ? "group-hover/ib:hidden group-focus-visible/ib:hidden" : ""}`} />
			{HoverIcon && (
				<HoverIcon className={`${iconSizes[size]} text-on-surface-variant hidden group-hover/ib:block group-focus-visible/ib:block`} />
			)}
		</Button>
	)
}

