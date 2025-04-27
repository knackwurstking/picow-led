import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
    build: {
        copyPublicDir: false,
        rollupOptions: {
            input: {
                api: resolve(__dirname, 'scripts/api.js'),
            },
            output: {
                dir: "public/js/",
                entryFileNames: "api.js",
            },
        },
    },
})
