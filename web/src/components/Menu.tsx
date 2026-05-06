import { Menu as MenuPremitive } from "@base-ui/react";
import { EllipsisIcon, type LucideIcon } from "lucide-react";
import type { ReactNode } from "react";

import { IconButton } from "./IconButton.tsx";

type MenuProps = {
  children: ReactNode;
};

export function Menu({ children }: MenuProps) {
  return (
    <MenuPremitive.Root>
      <MenuPremitive.Trigger render={<IconButton icon={EllipsisIcon} size="xs" />} />
      <MenuPremitive.Portal>
        <MenuPremitive.Positioner side="bottom" align="end">
          <MenuPremitive.Popup className="flex min-h-200 w-220 flex-col gap-1 rounded-xl bg-background px-4 py-2 shadow">
            {children}
          </MenuPremitive.Popup>
        </MenuPremitive.Positioner>
      </MenuPremitive.Portal>
    </MenuPremitive.Root>
  );
}

type MenuItemProps = Omit<MenuPremitive.Item.Props, "className"> & {
  icon: LucideIcon;
  label: string;
};

export function MenuItem({ icon: Icon, label, ...props }: MenuItemProps) {
  return (
    <MenuPremitive.Item
      render={<button type="button" />}
      nativeButton
      className="state-layer flex h-40 items-center gap-8 rounded-lg px-12"
      {...props}
    >
      <Icon className="size-20" />
      {label}
    </MenuPremitive.Item>
  );
}
