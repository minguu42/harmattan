import { Field, Input as InputBase } from "@base-ui/react";

type props = {
  label: string;
  valueMissingMessage?: string;
} & Omit<InputBase.Props, "children" | "className" | "onChange">;

export function Input({ label, valueMissingMessage, ...props }: props) {
  return (
    <Field.Root>
      <Field.Label className="text-sm font-medium text-on-background">{label}</Field.Label>
      <InputBase
        {...props}
        className="h-36 w-full rounded-md border border-input bg-transparent px-12 py-4 text-base text-on-background placeholder-placeholder"
      />
      {valueMissingMessage && (
        <Field.Error className="text-sm text-destructive" match="valueMissing">
          {valueMissingMessage}
        </Field.Error>
      )}
    </Field.Root>
  );
}
