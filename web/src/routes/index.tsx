import { Form } from "@base-ui/react";
import { createFileRoute } from "@tanstack/react-router";
import { PencilIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";

import { Button } from "../components/Button.tsx";
import { IconButton } from "../components/IconButton.tsx";
import { InputNew } from "../components/Input.tsx";
import { Select } from "../components/Select.tsx";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const [name, setName] = useState("");

  function alertName() {
    alert(name);
    setName("");
  }

  return (
    <div className="mx-auto flex max-w-720 flex-col items-center">
      <div className="flex items-center gap-12 p-12">
        <IconButton icon={PencilIcon} size="sm" />
        <IconButton icon={PencilIcon} size="md" />
        <IconButton icon={PencilIcon} hoverIcon={Trash2Icon} size="sm" />
        <IconButton icon={PencilIcon} hoverIcon={Trash2Icon} size="md" />
      </div>
      <div className="flex items-center gap-12 p-12">
        <Button label="作成" color="primary" />
        <Button label="有効" color="secondary" />
        <Button label="リセット" color="outlined" />
        <Button label="キャンセル" color="text" />
        <Button label="削除" color="destructive" />
      </div>
      <Select />
      <Form onFormSubmit={alertName}>
        <InputNew
          required
          label="プロジェクト名"
          placeholder="プロジェクト名"
          valueMissingMessage="プロジェクト名を入力してください"
          value={name}
          onValueChange={(v) => setName(v)}
        />
        <div className="h-16" />
        <Button label="あいうえお" type="submit" />
      </Form>
    </div>
  );
}
