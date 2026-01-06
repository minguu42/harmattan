import {Button as BaseButton} from "@base-ui/react"

type Props = {
	label: string
	type?: "submit" | "reset" | "button" | undefined
	onClick?: () => void
}

export function Button({label, type, onClick}: Props) {
	return (
		<BaseButton type={type} className="state-layer-parent rounded-xl text-primary h-12 min-w-12 text-sm" onClick={onClick}>
			<div className="state-layer bg-primary h-8" />
			<div className="state-layer-ring h-8" />
			{label}
		</BaseButton>
	)
}
