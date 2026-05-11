import { Avatar } from "@base-ui/react";
import { Link } from "@tanstack/react-router";
import {
  ArchiveIcon,
  FolderOpenDotIcon,
  InboxIcon,
  type LucideIcon,
  PencilIcon,
  PlusIcon,
  StarIcon,
  SunIcon,
  SunMoonIcon,
  TagIcon,
  Trash2Icon,
} from "lucide-react";
import { type Dispatch, type SetStateAction, useState } from "react";

import {
  type Project,
  useCreateProject,
  useDeleteProject,
  useProjects,
  useUpdateProject,
} from "../api/projects.ts";
import {type Tag, useCreateTag, useDeleteTag, useTags, useUpdateTag} from "../api/tags.ts"
import { Button } from "./Button.tsx";
import { Dialog, DialogTitle } from "./Dialog.tsx";
import { SelectField, TextField } from "./Field.tsx";
import { Form } from "./Form.tsx";
import { IconButton } from "./IconButton.tsx";
import { Menu, MenuItem } from "./Menu.tsx";

export function NavigationDrawer() {
  return (
    <div className="min-h-svh w-280 border-r border-border bg-sidebar">
      <div className="flex h-56 items-center gap-8 px-12">
        <Avatar.Root className="inline-grid size-32 place-items-center overflow-hidden rounded-full">
          <Avatar.Fallback className="grid size-full place-items-center bg-red-600 text-base text-white select-none">
            YF
          </Avatar.Fallback>
        </Avatar.Root>
        <div className="text-sm font-semibold text-on-sidebar">ボッスン</div>
      </div>
      <ul className="mt-8">
        <Indicator icon={SunIcon} label="今日" />
        <Indicator icon={SunMoonIcon} label="明日" />
        <Indicator icon={StarIcon} label="お気に入り" />
        <Indicator icon={InboxIcon} label="すべてのタスク" />
      </ul>
      <ProjectIndicatorList />
      <TagIndicatorList />
    </div>
  );
}

function Indicator({ icon: Icon, label }: { icon: LucideIcon; label: string }) {
  return (
    <li>
      <Link
        to="/"
        className="state-layer mx-8 flex h-36 items-center gap-8 rounded-lg px-8 text-sm text-on-sidebar"
      >
        <Icon className="size-24" />
        {label}
      </Link>
    </li>
  );
}

function ProjectIndicatorList() {
  const [open, setOpen] = useState(false);
  const { data: projects, error, isPending, isError } = useProjects();

  if (isPending) {
    return <span>Loading...</span>;
  }
  if (isError) {
    return <span>Error: {error.message}</span>;
  }
  return (
    <>
      <div className="mt-16 flex h-32 items-center justify-between pr-16 pl-12">
        <div className="text-sm font-semibold text-on-sidebar">プロジェクト</div>
        <IconButton icon={PlusIcon} size="sm" onClick={() => setOpen(true)} />
        <ProjectCreateDialog open={open} setOpen={setOpen} />
      </div>
      <ul>
        {projects.map((p) => (
          <ProjectIndicator key={p.id} project={p} />
        ))}
      </ul>
    </>
  );
}

const colors = [
  { label: "青色", value: "blue" },
  { label: "茶色", value: "brown" },
  { label: "白色", value: "default" },
  { label: "灰色", value: "gray" },
  { label: "緑色", value: "green" },
  { label: "橙色", value: "orange" },
  { label: "桃色", value: "pink" },
  { label: "紫色", value: "purple" },
  { label: "赤色", value: "red" },
  { label: "黄色", value: "yellow" },
];

type ProjectCreateDialogProps = {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

function ProjectCreateDialog({ open, setOpen }: ProjectCreateDialogProps) {
  const createProject = useCreateProject();

  function handleFormSubmit(formValues: Record<string, any>) {
    const name = formValues["name"].trim();
    const color = formValues["color"];
    if (name === "" || color === "") {
      return;
    }

    createProject.mutate(
      { name: name, color: color },
      {
        onSuccess: () => {
          setOpen(false);
        },
      },
    );
  }

  return (
    <Dialog open={open} setOpen={setOpen}>
      <DialogTitle>プロジェクト作成</DialogTitle>
      <Form className="flex flex-col gap-16" onFormSubmit={handleFormSubmit}>
        <TextField
          name="name"
          label="プロジェクト名"
          placeholder="プロジェクト名"
          required
          missingMessage="プロジェクト名は必須です"
        />
        <SelectField
          name="color"
          label="プロジェクトカラー"
          items={colors}
          required
          missingMessage="プロジェクトカラーは必須です"
        />
        <div className="flex gap-8">
          <div className="flex-1" />
          <Button label="キャンセル" color="text" onClick={() => setOpen(false)} />
          <Button type="submit" label="作成" />
        </div>
      </Form>
    </Dialog>
  );
}

type ProjectUpdateDialogProps = {
  project: Project;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

function ProjectUpdateDialog({ project, open, setOpen }: ProjectUpdateDialogProps) {
  const updateProject = useUpdateProject();

  function handleFormSubmit(formValues: Record<string, any>) {
    const name = formValues["name"].trim();
    const color = formValues["color"];
    if (name === "" || color === "") {
      return;
    }

    updateProject.mutate(
      { projectID: project.id, name, color },
      {
        onSuccess: () => {
          setOpen(false);
        },
      },
    );
  }

  return (
    <Dialog open={open} setOpen={setOpen}>
      <DialogTitle>プロジェクト変更</DialogTitle>
      <Form className="flex flex-col gap-16" onFormSubmit={handleFormSubmit}>
        <TextField
          name="name"
          label="プロジェクト名"
          defaultValue={project.name}
          required
          missingMessage="プロジェクト名は必須です"
        />
        <SelectField
          name="color"
          label="プロジェクトカラー"
          items={colors}
          defaultValue={project.color}
          required
          missingMessage="プロジェクトカラーは必須です"
        />
        <div className="flex gap-8">
          <div className="flex-1" />
          <Button label="キャンセル" color="text" onClick={() => setOpen(false)} />
          <Button type="submit" label="作成" />
        </div>
      </Form>
    </Dialog>
  );
}

function ProjectIndicator({ project }: { project: Project }) {
  const [open, setOpen] = useState(false);
  const deleteProject = useDeleteProject();

  return (
    <li className="group relative">
      <Link
        to="/projects/$projectID/tasks"
        params={{ projectID: project.id }}
        className="state-layer-with-sibling mx-8 flex h-36 items-center gap-8 rounded-lg px-8 text-sm text-on-background"
      >
        <FolderOpenDotIcon />
        {project.name}
      </Link>
      <div className="invisible absolute top-1/2 right-16 -translate-y-1/2 group-focus-within:visible group-hover:visible">
        <Menu>
          <MenuItem icon={PencilIcon} label="変更" onClick={() => setOpen(true)} />
          <MenuItem icon={ArchiveIcon} label="アーカイブ" />
          <MenuItem
            icon={Trash2Icon}
            label="削除"
            color="destructive"
            onClick={() => deleteProject.mutate(project.id)}
          />
        </Menu>
        <ProjectUpdateDialog project={project} open={open} setOpen={setOpen} />
      </div>
    </li>
  );
}

function TagIndicatorList() {
	const [open, setOpen] = useState(false);
	const { data: tags, error, isPending, isError } = useTags();

	if (isPending) {
		return <span>Loading...</span>;
	}
	if (isError) {
		return <span>Error: {error.message}</span>;
	}
	return (
		<>
			<div className="mt-16 flex h-32 items-center justify-between pr-16 pl-12">
				<div className="text-sm font-semibold text-on-sidebar">タグ</div>
				<IconButton icon={PlusIcon} size="sm" onClick={() => setOpen(true)} />
				<TagCreateDialog open={open} setOpen={setOpen} />
			</div>
			<ul>
				{tags.map((t) => (
					<TagIndicator key={t.id} tag={t} />
				))}
			</ul>
		</>
	);
}

function TagIndicator({ tag }: { tag: Tag }) {
	const [open, setOpen] = useState(false);
	const deleteTag = useDeleteTag();

	return (
		<li className="group relative">
			<Link
				to="/"
				className="state-layer-with-sibling mx-8 flex h-36 items-center gap-8 rounded-lg px-8 text-sm text-on-background"
			>
				<TagIcon />
				{tag.name}
			</Link>
			<div className="invisible flex absolute top-1/2 right-16 -translate-y-1/2 group-focus-within:visible group-hover:visible">
				<IconButton icon={PencilIcon} size="sm" onClick={() => setOpen(true)} />
				<IconButton icon={Trash2Icon} size="sm" onClick={() => deleteTag.mutate(tag.id)} />
			</div>
			<TagUpdateDialog tag={tag} open={open} setOpen={setOpen}/>
		</li>
	);
}

type DialogProps = {
	open: boolean;
	setOpen: Dispatch<SetStateAction<boolean>>;
};

function TagCreateDialog({ open, setOpen }: DialogProps) {
	const createTag = useCreateTag();

	function handleFormSubmit(values: Record<string, any>) {
		const name = values["name"].trim();
		if (name === "") {
			return;
		}

		createTag.mutate(
			{ name: name },
			{
				onSuccess: () => {
					setOpen(false);
				},
			},
		);
	}

	return (
		<Dialog open={open} setOpen={setOpen}>
			<DialogTitle>タグ作成</DialogTitle>
			<Form className="flex flex-col gap-16" onFormSubmit={handleFormSubmit}>
				<TextField
					name="name"
					label="タグ名"
					placeholder="タグ名"
					required
					missingMessage="タグ名は必須です"
				/>
				<div className="flex gap-8">
					<div className="flex-1" />
					<Button label="キャンセル" color="text" onClick={() => setOpen(false)} />
					<Button type="submit" label="作成" />
				</div>
			</Form>
		</Dialog>
	);
}

type TagUpdateDialogProps = {
	tag: Tag;
	open: boolean;
	setOpen: Dispatch<SetStateAction<boolean>>;
};

function TagUpdateDialog({ tag, open, setOpen }: TagUpdateDialogProps) {
	const updateTag = useUpdateTag();

	function handleFormSubmit(values: Record<string, any>) {
		const name = values["name"].trim();
		if (name === "") {
			return;
		}

		updateTag.mutate(
			{ tagID: tag.id, name },
			{
				onSuccess: () => {
					setOpen(false);
				},
			},
		);
	}

	return (
		<Dialog open={open} setOpen={setOpen}>
			<DialogTitle>タグ変更</DialogTitle>
			<Form className="flex flex-col gap-16" onFormSubmit={handleFormSubmit}>
				<TextField
					name="name"
					label="タグ名"
					defaultValue={tag.name}
					placeholder="タグ名"
					required
					missingMessage="タグ名は必須です"
				/>
				<div className="flex gap-8">
					<div className="flex-1" />
					<Button label="キャンセル" color="text" onClick={() => setOpen(false)} />
					<Button type="submit" label="変更" />
				</div>
			</Form>
		</Dialog>
	);
}
