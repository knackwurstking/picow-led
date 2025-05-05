(() => {
    const utils = require("../lib/utils");

    let timeout = null;
    window.addEventListener("focus", () => {
        if (timeout !== null) {
            return;
        }

        // NOTE: This focus event gets triggered twice after tab switching,
        //       this is just a workaround to prevent this
        timeout = setTimeout(async () => {
            try {
                const resp = await fetch(
                    // @ts-ignore
                    process.env.SERVER_PATH_PREFIX + "/api/ping",
                );
                const data = await resp.text();
                if (data === "pong") {
                    utils.setOnlineIndicatorState(true);
                } else {
                    utils.setOnlineIndicatorState(false);
                }
            } catch (err) {
                console.error(err);
                utils.setOnlineIndicatorState(false);
            }

            timeout = null;
        }, 100);
    });

    window.addEventListener("pageshow", () => {
        // Trigger the focus event once
        window.dispatchEvent(new Event("focus"));
    });
})();
