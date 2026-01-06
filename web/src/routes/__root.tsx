import {createRootRoute, Outlet} from "@tanstack/react-router"
import {NavigationDrawer} from "./-NavigationDrawer.tsx"

export const Route = createRootRoute({component: Root})

function Root() {
	return (
		<div className="isolate">
			<NavigationDrawer />
			<Outlet />
		</div>
	)
}
