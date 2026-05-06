import { createFileRoute } from "@tanstack/react-router";
import { PencilIcon, Trash2Icon } from "lucide-react";

import { Button } from "../components/Button.tsx";
import { IconButton } from "../components/IconButton.tsx";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="mx-auto flex w-600 items-center justify-center gap-4">
      <IconButton icon={PencilIcon} size="xs" />
      <IconButton icon={PencilIcon} size="sm" />
      <IconButton icon={PencilIcon} hoverIcon={Trash2Icon} size="sm" />
      <Button label="キャンセル" size="sm" color="text" />
      <Button label="送信" size="sm" color="filled" />
    </div>
  );
}
