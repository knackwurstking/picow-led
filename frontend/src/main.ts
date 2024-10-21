import "./app/picow-app";

import { registerSW } from "virtual:pwa-register";

registerSW({
    onRegistered(r) {
        console.debug("ServiceWorkerRegistration:", r);
        if (!r) return;

        setTimeout(async () => {
            try {
                console.debug(`PWA Update service...`);
                await r.update(); // For now do auto update all the time
            } catch (err) {
                console.warn(`PWA Auto update: ${err}`);
            }
        });
    },
});
