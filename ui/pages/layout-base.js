(() => {
    // TODO: Starting the WebSocket handler either in "focus event" or in "pageshow"

    let timeout = null;
    window.addEventListener("focus", () => {
        if (timeout !== null) {
            return;
        }

        // NOTE: This focus event gets triggered twice after tab switching,
        //       this is just a workaround to prevent running this twice
        timeout = setTimeout(async () => {
            try {
                const resp = await fetch(
                    // @ts-ignore
                    process.env.SERVER_PATH_PREFIX + "/api/ping",
                );
                const data = await resp.text();
                if (data === "pong") {
                    window.utils.setOnlineIndicatorState(true);
                } else {
                    window.utils.setOnlineIndicatorState(false);
                }
            } catch (err) {
                console.error(err);
                window.utils.setOnlineIndicatorState(false);
            }

            timeout = null;
        });
    });

    window.addEventListener("pageshow", () => {
        // Trigger the focus event once
        window.dispatchEvent(new Event("focus"));
    });
})();
