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
                api: resolve(__dirname, "ui/api.js"),
                ws: resolve(__dirname, "ui/ws.js"),
                utils: resolve(__dirname, "ui/utils.js"),
            },
            output: {
                dir: "public/js/",
                entryFileNames: "[name].js",
            },
        },
    },
});
