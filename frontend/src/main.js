import "../node_modules/ui/css/main.css";

import { register } from "ui";
import { registerSW } from "virtual:pwa-register";
import { createApp } from "./app";

registerSW({
    onRegistered(r) {
        console.debug("ServiceWorkerRegistration:", r);
        if (!r) return;

        setTimeout(async () => {
            try {
                console.debug(`PWA Update service...`);
                await r.update(); // NOTE: for now do auto update all the time
            } catch (err) {
                console.warn(`PWA Auto update: ${err}`);
            }
        });
    },
});

register();

const app = createApp();
document.querySelector(`div#app`).replaceWith(app.element);
