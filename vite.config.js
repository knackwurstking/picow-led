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
                main: "./script/main.ts",

                // Layouts
                "layouts/base": "./script/layouts/base.ts",

                // Pages
                "content/devices-address": "./script/content/address/main.ts",
                "content/devices": "./script/content/devices/main.ts",
                "content/settings": "./script/content/settings/main.ts",
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
