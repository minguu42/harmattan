import { Field, Form } from "@base-ui/react";
import { createFileRoute } from "@tanstack/react-router";
import { CheckIcon, CircleIcon, FolderOpenDotIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";

import { useProject } from "../api/projects.ts";
import { useCompleteTask, useCreateTask, useDeleteTask, useTasks } from "../api/tasks.ts";
import { IconButton } from "../components/IconButton.tsx";

export const Route = createFileRoute("/projects/$projectID/tasks")({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectID } = Route.useParams();
  return (
    <div className="mx-auto max-w-712">
      <div className="h-36" />
      <TaskListHeader projectID={projectID} />
      <div className="h-16" />
      <TaskAddField projectID={projectID} />
      <TaskList projectID={projectID} />
    </div>
  );
}

function TaskListHeader({ projectID }: { projectID: string }) {
  const { data: project, error, isPending, isError } = useProject(projectID);

  if (isPending) {
    return <span>Loading...</span>;
  }
  if (isError) {
    return <span>Error: {error.message}</span>;
  }
  return (
    <div className="flex items-center text-on-surface">
      <FolderOpenDotIcon size={36} />
      <h1 className="ml-4 text-[32px]">{project.name}</h1>
    </div>
  );
}

function TaskAddField({ projectID }: { projectID: string }) {
  const [name, setName] = useState("");
  const createTask = useCreateTask(projectID);

  function addTask() {
    if (name.trim() === "") return;
    createTask.mutate(name.trim(), {
      onSuccess: () => setName(""),
    });
  }

  return (
    <div className="flex h-40 items-center p-4">
      <PlusIcon className="text-[#8e8a90]" />
      <Form
        className="ml-8 w-full"
        onSubmit={(e) => {
          e.preventDefault();
          addTask();
        }}
      >
        <Field.Root>
          <Field.Control
            className="w-full p-4 text-on-surface placeholder:text-[#8e8a90]"
            required
            placeholder="タスク名"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </Field.Root>
      </Form>
    </div>
  );
}

function TaskList({ projectID }: { projectID: string }) {
  const { data: tasks, error, isPending, isError } = useTasks(projectID);
  const completeTask = useCompleteTask(projectID);
  const deleteTask = useDeleteTask(projectID);

  if (isPending) {
    return <span>Loading...</span>;
  }
  if (isError) {
    return <span>Error: {error.message}</span>;
  }
  return (
    <ul className="flex flex-col gap-px">
      {tasks.map((t) => (
        <li key={t.id} className="flex h-40 items-center gap-12">
          <IconButton
            icon={CircleIcon}
            hoverIcon={CheckIcon}
            onClick={() => completeTask.mutate(t.id)}
          />
          <p className="text-base text-on-surface">{t.name}</p>
          <div className="flex-1" />
          <IconButton icon={Trash2Icon} onClick={() => deleteTask.mutate(t.id)} />
        </li>
      ))}
    </ul>
  );
}
