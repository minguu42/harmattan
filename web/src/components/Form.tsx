import { Form as FormBase } from "@base-ui/react";

type props = Omit<FormBase.Props, "onSubmit">;

export function Form({ ...props }: props) {
  return <FormBase {...props} />;
}
