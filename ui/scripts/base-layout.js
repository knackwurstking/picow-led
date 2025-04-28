(() => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    window.addEventListener("focus", () => {
        fetch(__SERVER_PATH_PREFIX__ + "/api/ping")
            .then((r) => r.text())
            .then((d) => {
                if (d === "pong") {
                    utils.setOnlineIndicator(true);
                } else {
                    utils.setOnlineIndicator(false);
                }
            })
            .catch((err) => {
                console.error(err);
                utils.setOnlineIndicator(false);
            });
    });
})();
