import { Button as ButtonBase } from "@base-ui/react";
import { tv, type VariantProps } from "tailwind-variants";

const button = tv({
  base: "state-layer h-36 rounded-lg px-12 text-sm font-medium",
  variants: {
    color: {
      primary: "bg-primary text-on-primary",
      secondary: "bg-secondary text-on-secondary",
      outlined: "text-on-backgroud border border-border",
      text: "text-on-background",
      destructive: "bg-destructive-container text-destructive",
    },
  },
  defaultVariants: {
    color: "primary",
  },
});

type Props = Omit<ButtonBase.Props, "className" | "children"> & {
  label: string;
} & VariantProps<typeof button>;

export function Button({ label, color, ...props }: Props) {
  return (
    <ButtonBase className={button({ color: color })} {...props}>
      {label}
    </ButtonBase>
  );
}
