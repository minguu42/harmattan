import { Button } from "@base-ui/react";
import type { LucideIcon } from "lucide-react";
import { tv, type VariantProps } from "tailwind-variants";

const iconButton = tv({
  slots: {
    button: "group/ib state-layer grid place-items-center bg-transparent text-on-background",
    icon: "",
    hoverIcon: "",
  },
  variants: {
    size: {
      sm: {
        button: "size-32 rounded-lg",
        icon: "size-20",
        hoverIcon: "size-20",
      },
      md: {
        button: "size-40 rounded-xl",
        icon: "size-24",
        hoverIcon: "size-24",
      },
    },
    useHoverIcon: {
      true: {
        icon: "group-hover/ib:hidden group-focus-visible/ib:hidden",
        hoverIcon: "hidden group-hover/ib:block group-focus-visible/ib:block",
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
});

type Props = {
  icon: LucideIcon;
  hoverIcon?: LucideIcon;
} & Omit<Button.Props, "children" | "className"> &
  Omit<VariantProps<typeof iconButton>, "useHoverIcon">;

export function IconButton({ icon: Icon, hoverIcon: HoverIcon, size, ...props }: Props) {
  const { button, icon, hoverIcon } = iconButton({
    size: size,
    useHoverIcon: HoverIcon !== undefined,
  });
  return (
    <Button className={button()} {...props}>
      <Icon className={icon()} />
      {HoverIcon && <HoverIcon className={hoverIcon()} />}
    </Button>
  );
}
