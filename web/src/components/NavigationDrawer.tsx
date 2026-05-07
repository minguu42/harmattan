import { Avatar, Dialog, Field, Form } from "@base-ui/react";
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
import { useState } from "react";

import { type Project } from "../api/projects.ts";
import { useCreateProject, useDeleteProject, useProjects } from "../api/projects.ts";
import { Button } from "./Button.tsx";
import { IconButton } from "./IconButton.tsx";
import { Input } from "./Input.tsx";
import { Menu, MenuItem } from "./Menu.tsx";

export function NavigationDrawer() {
  return (
    <div className="min-h-svh w-280 border-r border-border bg-sidebar">
      <div className="flex h-56 items-center gap-8 px-12">
        <Avatar.Root className="inline-flex size-32 items-center justify-center overflow-hidden rounded-full bg-red-600 align-middle text-base text-white select-none">
          YF
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
        <ProjectCreateButtonWithDialog />
      </div>
      <ul>
        {projects.map((p) => (
          <ProjectIndicator key={p.id} project={p} />
        ))}
      </ul>
    </>
  );
}

function ProjectCreateButtonWithDialog() {
  const [open, setOpen] = useState(false);

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <IconButton icon={PlusIcon} size="sm" onClick={() => setOpen(true)} />
      <Dialog.Portal>
        <Dialog.Backdrop forceRender className="fixed inset-0 bg-scrim" />
        <ProjectCreateDialog closePopup={() => setOpen(false)} />
      </Dialog.Portal>
    </Dialog.Root>
  );
}

function ProjectCreateDialog({ closePopup }: { closePopup: () => void }) {
  const [name, setName] = useState("");
  const createProject = useCreateProject();

  function addProject() {
    if (name.trim() === "") {
      return;
    }

    createProject.mutate(name.trim(), {
      onSuccess: () => {
        setName("");
        closePopup();
      },
    });
  }

  return (
    <Dialog.Popup className="fixed top-1/2 left-1/2 -mt-32 w-400 -translate-x-1/2 -translate-y-1/2 rounded-xl border border-border bg-background p-16">
      <Dialog.Title className="text-on-surface text-base font-medium">
        プロジェクト作成
      </Dialog.Title>
      <Form
        onSubmit={(e) => {
          e.preventDefault();
          addProject();
        }}
      >
        <Field.Root className="my-16 flex flex-col gap-4">
          <Field.Label className="text-sm font-medium text-on-background">
            プロジェクト名
          </Field.Label>
          <Input
            type="text"
            required
            placeholder="プロジェクト名"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </Field.Root>
        <div className="flex">
          <div className="flex-1" />
          <Button type="submit" label="作成" />
        </div>
      </Form>
    </Dialog.Popup>
  );
}

function ProjectIndicator({ project }: { project: Project }) {
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
          <MenuItem icon={PencilIcon} label="変更" />
          <MenuItem icon={ArchiveIcon} label="アーカイブ" />
          <MenuItem
            icon={Trash2Icon}
            label="削除"
            onClick={() => deleteProject.mutate(project.id)}
          />
        </Menu>
      </div>
    </li>
  );
}
