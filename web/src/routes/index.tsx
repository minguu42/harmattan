import { createFileRoute } from "@tanstack/react-router";
import { PencilIcon, Trash2Icon } from "lucide-react";

import { Button } from "../components/Button.tsx";
import { SelectField, TextField } from "../components/Field.tsx";
import { Form } from "../components/Form.tsx";
import { IconButton } from "../components/IconButton.tsx";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
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

      <Form
        className="flex flex-col gap-16"
        onFormSubmit={(formValues) => {
          const name = formValues["name"];
          const color = formValues["color"];
          alert(`name: ${name}, color: ${color}`);
        }}
      >
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
          placeholder="プロジェクトカラー"
          required
          missingMessage="プロジェクトカラーは必須です"
        />
        <Button type="submit" label="送信" />
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
