import { Button as ButtonBase } from "@base-ui/react";
import { tv, type VariantProps } from "tailwind-variants";

const button = tv({
  base: "state-layer font-medium",
  variants: {
    size: {
      xs: "h-32 rounded-lg px-12 text-sm",
      sm: "h-40 rounded-lg px-16 text-sm",
      md: "h-56 rounded-xl px-24",
    },
    color: {
      filled: "bg-primary text-white",
      text: "text-primary",
    },
  },
  defaultVariants: {
    size: "sm",
    color: "filled",
  },
});

type Props = Omit<ButtonBase.Props, "className" | "children"> & {
  label: string;
} & VariantProps<typeof button>;

export function Button({ label, size, color, ...props }: Props) {
  return (
    <ButtonBase className={button({ size: size, color: color })} {...props}>
      {label}
    </ButtonBase>
  );
}
