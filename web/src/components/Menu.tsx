import { Menu as MenuBase } from "@base-ui/react";
import { EllipsisIcon, type LucideIcon } from "lucide-react";
import type { ReactNode } from "react";

import { IconButton } from "./IconButton.tsx";

type MenuProps = {
  children: ReactNode;
};

export function Menu({ children }: MenuProps) {
  return (
    <MenuBase.Root>
      <MenuBase.Trigger render={<IconButton icon={EllipsisIcon} size="sm" />} />
      <MenuBase.Portal>
        <MenuBase.Positioner side="bottom" align="end">
          <MenuBase.Popup className="flex w-full min-w-0 flex-col gap-1 rounded-xl bg-background p-4 shadow">
            {children}
          </MenuBase.Popup>
        </MenuBase.Positioner>
      </MenuBase.Portal>
    </MenuBase.Root>
  );
}

type MenuItemProps = Omit<MenuBase.Item.Props, "className"> & {
  icon: LucideIcon;
  label: string;
};

export function MenuItem({ icon: Icon, label, ...props }: MenuItemProps) {
  return (
    <MenuBase.Item
      className="state-layer flex h-32 min-w-160 items-center gap-8 rounded-lg px-8 py-4 text-sm text-on-background"
      {...props}
    >
      <Icon className="size-20" />
      {label}
    </MenuBase.Item>
  );
}
