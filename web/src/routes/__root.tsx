import { createRootRoute, Outlet } from "@tanstack/react-router";
import { PanelLeftIcon } from "lucide-react";
import { useState } from "react";

import { IconButton } from "../components/IconButton.tsx";
import { NavigationDrawer } from "../components/NavigationDrawer.tsx";

export const Route = createRootRoute({ component: Root });

function Root() {
  const [open, setOpen] = useState(false);
  return (
    <div className="isolate flex min-h-svh bg-background">
      {open && <NavigationDrawer />}
      <div className="flex-1">
        <Header toggleDrawer={() => setOpen((prev) => !prev)} />
        <Outlet />
      </div>
    </div>
  );
}

function Header({ toggleDrawer }: { toggleDrawer: () => void }) {
  return (
    <header className="flex h-56 items-center px-8">
      <IconButton icon={PanelLeftIcon} onClick={toggleDrawer} />
    </header>
  );
}
