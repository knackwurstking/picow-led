import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

const icons = [
    {
        src: "/icons/pwa-64x64.png",
        sizes: "64x64",
        type: "image/png",
    },
    {
        src: "/icons/pwa-192x192.png",
        sizes: "192x192",
        type: "image/png",
    },
    {
        src: "/icons/pwa-512x512.png",
        sizes: "512x512",
        type: "image/png",
    },
    {
        src: "/icons/maskable-icon-512x512.png",
        sizes: "512x512",
        type: "image/png",
        purpose: "maskable",
    },
];

const screenshots = [];

/** @type {import("vite-plugin-pwa").VitePWAOptions} */
const manifest = {
    strategies: "generateSW",
    registerType: "autoUpdate",
    includeAssets: [
        "/fonts/Recursive_VF_1.085--subset_range_english_basic.woff2",
        "/fonts/Recursive_VF_1.085--subset_range_latin_1_punc.woff2",
        "/fonts/Recursive_VF_1.085--subset_range_remaining.woff2",
        "/fonts/Recursive_VF_1.085--subset_range_latin_1.woff2",
        "/fonts/Recursive_VF_1.085--subset_range_latin_ext.woff2",
        "/fonts/Recursive_VF_1.085--subset_range_vietnamese.woff2",
    ],
    manifest: {
        name: "PicoW LED",
        short_name: "picow-led",
        icons: icons,
        screenshots: screenshots,
        theme_color: "#09090b",
        background_color: "#09090b",
        display: "standalone",
        scope: ".",
        start_url: "./",
        orientation: "any",
        //publicPath: "/picow-led.github.io",
    },
};

export default defineConfig({
    //server: {
    //    proxy: {
    //        "/api": "http://localhost:50833",
    //        "/events/device": "ws://localhost:50833",
    //    },
    //},
    clearScreen: false,
    plugins: [VitePWA(manifest)],
});
