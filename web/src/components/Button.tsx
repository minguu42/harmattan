import { Button as BButton } from "@base-ui/react"

type Props = {
	label: string
	style?: "filled" | "outlined" | "text"
	type?: "submit" | "reset" | "button"
	onClick?: () => void
}

const containerColors = {
	"filled": "bg-primary",
	"outlined": "outline-1 outline-outline",
	"text": "",
}

const layerColors = {
	"filled": "bg-on-primary",
	"outlined": "bg-on-surface-variant",
	"text": "bg-primary",
}

const labelColors = {
	"filled": "text-on-primary",
	"outlined": "text-on-surface-variant",
	"text": "text-primary",
}

export function Button({ label, style = "text", type = "button", onClick }: Props) {
	return (
		<BButton type={type} className="group/b h-12 min-w-12 flex items-center focus-visible:outline-none" onClick={onClick}>
			<div className={`relative h-10 px-4 rounded-xl flex items-center ${containerColors[style]}`}>
				<div className={`absolute inset-0 rounded-[inherit] opacity-0 group-hover/b:opacity-10 group-focus-visible/b:opacity-10 group-active/b:opacity-12 ${layerColors[style]}`} />
				<div className="absolute inset-0 rounded-[inherit] group-focus-visible/b:outline-2 group-focus-visible/b:outline-focus-ring" />
				<span className={`text-sm font-medium ${labelColors[style]}`}>{label}</span>
			</div>
		</BButton>
	)
}
