import {StrictMode} from "react"
import {createRoot} from "react-dom/client"
import {createRouter, RouterProvider} from "@tanstack/react-router"
import {routeTree} from "./routeTree.gen.ts"
import "./index.css"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query"

const client = new QueryClient()
const router = createRouter({routeTree})

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router
	}
}

const rootElement = document.getElementById("root")
if (rootElement && !rootElement.innerHTML) {
	createRoot(rootElement).render(
		<StrictMode>
			<QueryClientProvider client={client}>
				<RouterProvider router={router} />
			</QueryClientProvider>
		</StrictMode>,
	)
}
