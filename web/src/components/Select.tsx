import { Select as SelectBase } from "@base-ui/react";
import { CheckIcon, ChevronsUpDownIcon } from "lucide-react";

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

type props = {
  label: string;
};

export function Select({ label }: props) {
  return (
    <SelectBase.Root items={colors}>
      <div>
        <SelectBase.Label className="text-sm font-medium text-on-background">
          {label}
        </SelectBase.Label>
        <SelectBase.Trigger className="flex h-36 w-full items-center gap-4 rounded-md border border-input px-8 text-base">
          <SelectBase.Value
            className="data-placeholder:text-placeholder"
            placeholder="プロジェクトカラー"
          />
          <SelectBase.Icon className="text-placeholder">
            <ChevronsUpDownIcon />
          </SelectBase.Icon>
        </SelectBase.Trigger>
      </div>
      <SelectBase.Portal>
        <SelectBase.Positioner>
          <SelectBase.Popup className="min-w-(--anchor-width) bg-background p-4 shadow">
            <SelectBase.List>
              {colors.map(({ label, value }) => (
                <SelectBase.Item
                  key={label}
                  value={value}
                  className="state-layer flex h-32 items-center gap-4 rounded-lg px-4"
                >
                  <SelectBase.ItemIndicator>
                    <CheckIcon />
                  </SelectBase.ItemIndicator>
                  <SelectBase.ItemText>{label}</SelectBase.ItemText>
                </SelectBase.Item>
              ))}
            </SelectBase.List>
          </SelectBase.Popup>
        </SelectBase.Positioner>
      </SelectBase.Portal>
    </SelectBase.Root>
  );
}
