import * as path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import sunmaoFsVitePlugin from "@sunmao-ui/vite-plugin-fs";
import routes, { type RouteConfig } from "./src/routes";

const routeConfigs = routes.filter((route) => "name" in route) as RouteConfig[];

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    sunmaoFsVitePlugin({
      schemas: routeConfigs.map((route) => ({
        name: route.name,
        path: path.resolve(__dirname, `./src/application/${route.name}.json`),
      })),
      modulesDir: path.resolve(__dirname, "./src/modules"),
    }),
    react(),
  ],
  define: {
    // react-codemirror2 needs this
    global: "globalThis",
    // https://github.com/vitejs/vite/issues/1973
    "process.env": {},
    "process.platform": '"web"',
  },
  build: {
    rollupOptions: {
      input: {
        index: path.resolve(__dirname, "./index.html"),
        editor: path.resolve(__dirname, "./editor.html"),
      },
    },
  },
});
