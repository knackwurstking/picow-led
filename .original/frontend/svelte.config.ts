import adapter from "@sveltejs/adapter-static";
import { Config } from "@sveltejs/kit";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

const config: Config = {
    preprocess: vitePreprocess(),
    kit: {
        adapter: adapter(),
    },
};

export default config;
