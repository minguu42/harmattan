import { Field, Input as InputBase } from "@base-ui/react";

type Props = {
  label: string;
  valueMissingMessage?: string;
} & Omit<InputBase.Props, "className">;

export function Input({ label, valueMissingMessage, ...props }: Props) {
  return (
    <Field.Root>
      <Field.Label className="text-sm font-medium text-on-background">{label}</Field.Label>
      <InputBase
        {...props}
        className="h-36 w-full min-w-0 rounded-md border border-input bg-transparent px-12 py-4 text-base placeholder-placeholder"
      />
      {valueMissingMessage && (
        <Field.Error className="text-sm text-destructive" match="valueMissing">
          {valueMissingMessage}
        </Field.Error>
      )}
    </Field.Root>
  );
}
