import { defineConfig } from "vite";
import { babel } from "@rollup/plugin-babel";
import { viteRequire } from "vite-require";

export default defineConfig({
    plugins: [
        babel({
            babelHelpers: "bundled",
            presets: [
                [
                    "@babel/preset-env",
                    {
                        targets: {
                            browsers: ["chrome >= 55"],
                        },
                    },
                ],
            ],
        }),
        viteRequire(),
    ],

    build: {
        minify: false,
        copyPublicDir: false,
        emptyOutDir: false,
        rollupOptions: {
            input: {
                // Main, window
                main: "./ui/pages/main.js",

                // Layouts
                "layout-base": "./ui/pages/layout-base.js",

                // Pages
                "devices-address": "./ui/pages/devices/address/main.js",
                devices: "./ui/pages/devices/main.js",
                settings: "./ui/pages/settings/main.js",
            },
            output: {
                dir: "./public/js/",
                entryFileNames: "[name].js",
            },
        },
    },

    //esbuild: {
    //    legalComments: "inline",
    //},

    define: {
        "process.env.SERVER_PATH_PREFIX": JSON.stringify(
            process.env.SERVER_PATH_PREFIX,
        ),
    },
});
