import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
    build: {
        minify: false,
        copyPublicDir: false,
        rollupOptions: {
            input: {
                window: resolve(__dirname, "ui/window.js"),
                "layouts/base": resolve(__dirname, "ui/layouts/base.js"),
                "pages/devices": resolve(__dirname, "ui/pages/devices.js"),
                "pages/devices-addr": resolve(
                    __dirname,
                    "ui/pages/devices-addr.js",
                ),
                "pages/settings": resolve(__dirname, "ui/pages/settings.js"),
            },
            output: {
                dir: "templates/js/",
                entryFileNames: "[name].js",
            },
        },
    },

    esbuild: {
        legalComments: "inline",
    },
});
