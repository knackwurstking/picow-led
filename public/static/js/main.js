function registerServiceWorkers() {
    return new Promise(() => {
        // Check if the browser supports service workers, otherwise abort.
        if (!("serviceWorker" in navigator)) {
            console.warn("Browser doesn't support service workers");
            return;
        }

        window.addEventListener("load", function () {
            navigator.serviceWorker
                .register("service-worker.js")
                .then(function (reg) {
                    console.info("Service worker registered", reg);
                })
                .catch(function (err) {
                    console.error("Service worker registration failed:", err);
                });
        });
    });
}

registerServiceWorkers();
