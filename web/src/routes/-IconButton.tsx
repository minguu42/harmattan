import type {LucideIcon} from "lucide-react"
import {Button} from "@base-ui/react"

type Props = {
	icon: LucideIcon
	hoverIcon?: LucideIcon
	onClick?: () => void
}

export function IconButton({icon: Icon, hoverIcon: HoverIcon, onClick}: Props) {
	return (
		<Button className="state-layer-parent group/ib size-12 grid place-items-center rounded-xl" onClick={onClick}>
			<div className="state-layer size-10 bg-on-surface-variant" />
			<div className="state-layer-ring size-10" />
			<Icon className={`size-6 text-on-surface-variant ${HoverIcon ? "group-hover/ib:hidden group-focus-visible/ib:hidden" : ""}`} />
			{HoverIcon &&
          <HoverIcon className="hidden size-6 text-on-surface-variant group-hover/ib:block group-focus-visible/ib:block" />
			}
		</Button>
	)
}
