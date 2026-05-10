import { Avatar } from "@base-ui/react";
import { Link } from "@tanstack/react-router";
import {
  SunIcon,
  SunMoonIcon,
  StarIcon,
  InboxIcon,
  PlusIcon,
  TagIcon,
  type LucideIcon,
  FolderOpenDotIcon,
  PencilIcon,
  Trash2Icon,
  ArchiveIcon,
} from "lucide-react";
import { type Dispatch, type SetStateAction, useState } from "react";

import { type Project, useUpdateProject } from "../api/projects.ts";
import { useCreateProject, useDeleteProject, useProjects } from "../api/projects.ts";
import { Button } from "./Button.tsx";
import { Dialog, DialogTitle } from "./Dialog.tsx";
import { Form } from "./Form.tsx";
import { IconButton } from "./IconButton.tsx";
import { Input } from "./Input.tsx";
import { Menu, MenuItem } from "./Menu.tsx";
import { Select } from "./Select.tsx";

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
      <div className="mt-16 flex h-32 items-center justify-between pr-16 pl-12">
        <div className="text-sm font-semibold text-on-sidebar">タグ</div>
        <IconButton icon={PlusIcon} size="sm" />
      </div>
      <ul>
        <li className="flex h-36 items-center gap-8 px-16">
          <TagIcon />
          <div className="text-sm text-on-sidebar">タグ1</div>
        </li>
      </ul>
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
  const [name, setName] = useState("");
  const [color, setColor] = useState<string | null>(null);
  const createProject = useCreateProject();

  function handleFormSubmit() {
    if (name.trim() === "" || color == null) {
      return;
    }

    createProject.mutate(
      { name: name.trim(), color: color },
      {
        onSuccess: () => {
          setName("");
          setColor(null);
          setOpen(false);
        },
      },
    );
  }

  return (
    <Dialog open={open} setOpen={setOpen}>
      <DialogTitle>プロジェクト作成</DialogTitle>
      <Form className="flex flex-col gap-16" onFormSubmit={handleFormSubmit}>
        <Input
          required
          label="プロジェクト名"
          placeholder="プロジェクト名"
          value={name}
          onValueChange={(v) => setName(v)}
          valueMissingMessage="プロジェクト名は必須です"
        />
        <Select
          label="プロジェクトカラー"
          required
          valueMissingMessage="プロジェクトカラーは必須です"
          items={colors}
          value={color}
          onValueChange={(v) => setColor(v)}
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
  const [name, setName] = useState(project.name);
  const [color, setColor] = useState(project.color);
  const updateProject = useUpdateProject();

  function handleFormSubmit() {
    if (name.trim() === "" || color == null) {
      return;
    }

    updateProject.mutate(
      { projectID: project.id, name: name.trim(), color },
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
        <Input
          required
          label="プロジェクト名"
          placeholder="プロジェクト名"
          value={name}
          onValueChange={(v) => setName(v)}
          valueMissingMessage="プロジェクト名は必須です"
        />
        <Select
          label="プロジェクトカラー"
          required
          valueMissingMessage="プロジェクトカラーは必須です"
          items={colors}
          value={color}
          onValueChange={(v) => setColor(v)}
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
        className="text-on-surface state-layer-with-sibling mx-8 flex h-36 items-center gap-8 rounded-lg px-8 text-sm"
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
