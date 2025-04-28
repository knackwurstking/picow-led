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
                api: resolve(__dirname, "ui/scripts/api.js"),
                ws: resolve(__dirname, "ui/scripts/ws.js"),
                utils: resolve(__dirname, "ui/scripts/utils.js"),
                "service-worker": resolve(
                    __dirname,
                    "ui/scripts/service-worker.js",
                ),
            },
            output: {
                dir: "public/js/",
                entryFileNames: "[name].js",
            },
        },
    },
});
