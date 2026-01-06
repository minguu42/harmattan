import {defineConfig} from "vite"
import react from "@vitejs/plugin-react"
import {tanstackRouter} from "@tanstack/router-plugin/vite"
import tailwindcss from "@tailwindcss/vite"

// https://vite.dev/config/
export default defineConfig({
	root: "web",
	plugins: [
		// '@tanstack/router-plugin'は'@vitejs/plugin-react'より先に渡す必要がある
		tanstackRouter({
			target: "react",
			autoCodeSplitting: true,
		}),
		react(),
		tailwindcss(),
	],
	build: {
		outDir: "../dist",
	},
})
