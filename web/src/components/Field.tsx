import { Field, Input, Select } from "@base-ui/react";
import { CheckIcon, ChevronsUpDownIcon } from "lucide-react";

type TextFieldProps = {
  name: string;
  label: string;
  defaultValue?: string;
  placeholder?: string;
  required?: boolean;
  missingMessage?: string;
};

export function TextField({
  name,
  label,
  defaultValue,
  placeholder,
  required,
  missingMessage,
}: TextFieldProps) {
  return (
    <Field.Root name={name}>
      <Field.Label className="text-sm font-medium text-on-background">{label}</Field.Label>
      <Input
        defaultValue={defaultValue}
        placeholder={placeholder}
        required={required}
        className="h-36 w-full rounded-md border border-input bg-transparent px-12 py-4 text-base text-on-background placeholder-placeholder"
      />
      {missingMessage && (
        <Field.Error className="text-sm text-destructive">{missingMessage}</Field.Error>
      )}
    </Field.Root>
  );
}

type SelectFieldProps = {
  name: string;
  label: string;
  items: ReadonlyArray<{ value: any; label: string }>;
  defaultValue?: string;
  placeholder?: string;
  required?: boolean;
  missingMessage?: string;
};

export function SelectField({
  name,
  label,
  items,
  defaultValue,
  placeholder,
  required,
  missingMessage,
}: SelectFieldProps) {
  return (
    <Field.Root name={name}>
      <Select.Root items={items} required={required} defaultValue={defaultValue}>
        <Select.Label className="text-sm font-medium text-on-background">{label}</Select.Label>
        <Select.Trigger className="flex h-36 w-full items-center justify-between gap-4 rounded-md border border-input px-8">
          <Select.Value
            className="text-base data-placeholder:text-placeholder"
            placeholder={placeholder}
          />
          <Select.Icon className="text-placeholder">
            <ChevronsUpDownIcon />
          </Select.Icon>
        </Select.Trigger>
        <Select.Portal>
          <Select.Positioner>
            <Select.Popup className="min-w-(--anchor-width) bg-background p-4 shadow">
              <Select.List>
                {items.map(({ label, value }) => (
                  <Select.Item
                    key={label}
                    value={value}
                    className="state-layer flex h-32 items-center gap-4 rounded-lg px-4"
                  >
                    <Select.ItemIndicator>
                      <CheckIcon />
                    </Select.ItemIndicator>
                    <Select.ItemText>{label}</Select.ItemText>
                  </Select.Item>
                ))}
              </Select.List>
            </Select.Popup>
          </Select.Positioner>
        </Select.Portal>
      </Select.Root>
      {missingMessage && (
        <Field.Error className="text-sm text-destructive">{missingMessage}</Field.Error>
      )}
    </Field.Root>
  );
}
