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
                "base-layout": resolve(__dirname, "ui/scripts/base-layout.js"),
            },
            output: {
                dir: "public/js/",
                entryFileNames: "[name].js",
            },
        },
    },
    define: {
        __SERVER_PATH_PREFIX__: JSON.stringify(
            process.env.SERVER_PATH_PREFIX || "",
        ),
    },
});
