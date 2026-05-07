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

export function Select({}) {
  return (
    <div className="flex flex-col gap-1">
      <SelectBase.Root items={colors}>
        <SelectBase.Label className="text-sm text-on-background">
          プロジェクトカラー
        </SelectBase.Label>
        <SelectBase.Trigger className="flex h-36 min-w-160 items-center gap-3 border border-input px-8 text-base">
          <SelectBase.Value
            className="data-placeholder:text-placeholder"
            placeholder="プロジェクトカラー"
          />
          <SelectBase.Icon className="text-placeholder">
            <ChevronsUpDownIcon />
          </SelectBase.Icon>
        </SelectBase.Trigger>
        <SelectBase.Portal>
          <SelectBase.Positioner>
            <SelectBase.Popup>
              <SelectBase.ScrollUpArrow />
              <SelectBase.List>
                {colors.map(({ label, value }) => (
                  <SelectBase.Item key={label} value={value}>
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
    </div>
  );
}
