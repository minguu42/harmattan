import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite-plus";

// https://vite.dev/config/
export default defineConfig({
  fmt: {
    ignorePatterns: [],
    sortImports: true,
    sortTailwindcss: {
      functions: ["tv"],
      stylesheet: "./src/index.css",
    },
  },
  plugins: [
    // '@tanstack/router-plugin'は'@vitejs/plugin-react'より先に渡す必要がある
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
    tailwindcss(),
  ],
});
