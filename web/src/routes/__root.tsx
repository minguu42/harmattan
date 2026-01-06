import {createRootRoute, Outlet} from "@tanstack/react-router"
import {NavigationDrawer} from "./-NavigationDrawer.tsx"

export const Route = createRootRoute({component: RouteComponent})

function RouteComponent() {
	return (
		<div className="isolate min-h-screen bg-surface">
			<NavigationDrawer />
			<Outlet />
		</div>
	)
}
