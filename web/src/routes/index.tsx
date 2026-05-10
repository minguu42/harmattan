import { createFileRoute } from "@tanstack/react-router";
import { PencilIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";

import { Button } from "../components/Button.tsx";
import { Form } from "../components/Form.tsx";
import { IconButton } from "../components/IconButton.tsx";
import { Input } from "../components/Input.tsx";
import { Select } from "../components/Select.tsx";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const [name, setName] = useState("");
  const [color, setColor] = useState<string | null>(null);

  function alertName() {
    alert(`name: ${name}, color: ${color}`);
    setName("");
    setColor(null);
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
      <Form className="flex flex-col gap-16" onFormSubmit={alertName}>
        <Input
          required
          label="プロジェクト名"
          placeholder="プロジェクト名"
          valueMissingMessage="プロジェクト名を入力してください"
          value={name}
          onValueChange={(v) => setName(v)}
        />
        <Select
          required
					items={colors}
          value={color}
          onValueChange={(v) => setColor(v)}
          label="プロジェクトカラー"
        />
        <Button label="あいうえお" type="submit" />
      </Form>
    </div>
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
