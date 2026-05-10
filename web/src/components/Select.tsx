import { Field, Select as SelectBase } from "@base-ui/react";
import { CheckIcon, ChevronsUpDownIcon } from "lucide-react";

type props = {
  label: string;
  items: ReadonlyArray<{ label: string; value: any }>;
  valueMissingMessage?: string;
} & Omit<SelectBase.Root.Props<any>, "items">;

export function Select({ label, items, valueMissingMessage, ...props }: props) {
  return (
    <Field.Root>
      <SelectBase.Root items={items} {...props}>
        <SelectBase.Label className="text-sm font-medium text-on-background">
          {label}
        </SelectBase.Label>
        <SelectBase.Trigger className="flex h-36 w-full items-center justify-between gap-4 rounded-md border border-input px-8 text-base">
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
            <SelectBase.Popup className="min-w-(--anchor-width) bg-background p-4 shadow">
              <SelectBase.List>
                {items.map(({ label, value }) => (
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
      {valueMissingMessage && (
        <Field.Error match="valueMissing" className="text-sm text-destructive">
          {valueMissingMessage}
        </Field.Error>
      )}
    </Field.Root>
  );
}
