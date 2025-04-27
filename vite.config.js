import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
    build: {
        target: "es6",
        minify: false,
        copyPublicDir: false,
        rollupOptions: {
            input: {
                // TODO: Add server path prefix via environment here somehow
                api: resolve(__dirname, "ui/scripts/api.js"),
                ws: resolve(__dirname, "ui/scripts/ws.js"),
                utils: resolve(__dirname, "ui/scripts/utils.js"),
            },
            output: {
                dir: "public/js/",
                entryFileNames: "[name].js",
            },
        },
    },
});
