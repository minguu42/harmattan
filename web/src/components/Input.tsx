import { Input as InputBase } from "@base-ui/react";

export function Input({ ...props }: Omit<InputBase.Props, "className">) {
  return (
    <InputBase
      className="h-36 w-full min-w-0 rounded-md border border-input bg-transparent px-12 py-4 text-base placeholder-placeholder"
      {...props}
    />
  );
}
