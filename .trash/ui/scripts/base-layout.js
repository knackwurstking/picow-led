(() => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    window.addEventListener("focus", async () => {
        try {
            const resp = await fetch(__SERVER_PATH_PREFIX__ + "/api/ping");
            const data = await resp.text();
            if (data === "pong") {
                utils.setOnlineIndicator(true);
            } else {
                utils.setOnlineIndicator(false);
            }
        } catch (err) {
            console.error(err);
            utils.setOnlineIndicator(false);
        }
    });
})();
