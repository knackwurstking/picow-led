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
                main: "./ui/main.js",
                layout: "./ui/layout.js",
                settings: "./ui/pages/settings.js",
                devices: "./ui/pages/devices.js",
                "devices-address": "./ui/pages/devices-address.js",
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
