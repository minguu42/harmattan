import {createFileRoute} from "@tanstack/react-router"
import {MinusIcon, PencilIcon, PlusIcon, SettingsIcon, Trash2Icon} from "lucide-react"
import {Button} from "../components/Button.tsx"
import {IconButton} from "../components/IconButton.tsx"
import {Menu, MenuItem} from "../components/Menu.tsx"

export const Route = createFileRoute("/components")({
	component: RouteComponent,
})

function RouteComponent() {
	return (
		<div className="min-h-screen bg-surface p-8">
			<h1 className="text-3xl text-on-surface mb-8">コンポーネント一覧</h1>
			<div className="flex flex-col gap-12">
				<Section title="Button">
					<div className="flex items-center gap-4">
						<Button label="ボタン" />
						<Button label="作成" />
						<Button label="保存" />
					</div>
				</Section>
				<Section title="IconButton">
					<div className="flex items-center gap-4">
						<IconButton icon={PlusIcon} hoverIcon={MinusIcon} />
						<IconButton icon={PencilIcon} />
						<IconButton icon={Trash2Icon} />
						<IconButton icon={SettingsIcon} />
					</div>
				</Section>
				<Section title="Menu">
					<div className="flex items-center gap-8">
						<div className="flex flex-col items-center gap-2">
							<span className="text-sm text-on-surface-variant">基本</span>
							<Menu>
								<MenuItem icon={PencilIcon} label="変更" />
								<MenuItem icon={Trash2Icon} label="削除" />
							</Menu>
						</div>
						<div className="flex flex-col items-center gap-2">
							<span className="text-sm text-on-surface-variant">複数項目</span>
							<Menu>
								<MenuItem icon={PlusIcon} label="追加" />
								<MenuItem icon={PencilIcon} label="変更" />
								<MenuItem icon={SettingsIcon} label="設定" />
								<MenuItem icon={Trash2Icon} label="削除" />
							</Menu>
						</div>
					</div>
				</Section>
			</div>
		</div>
	)
}

type SectionProps = {
	title: string
	children: React.ReactNode
}

function Section({title, children}: SectionProps) {
	return (
		<section className="flex flex-col gap-4">
			<h2 className="text-xl text-on-surface border-b border-on-surface-variant pb-2">{title}</h2>
			<div className="pl-4">{children}</div>
		</section>
	)
}
