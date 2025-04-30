//{{ define "script-base-layout" }}
(() => {
    /** @type {number | null} */
    let timeout = null;
    window.addEventListener("focus", () => {
        if (timeout !== null) {
            return;
        }

        // NOTE: This focus event gets triggered twice after tab switching,
        //       this is just a workaround to prevent this
        timeout = setTimeout(async () => {
            /** @type {Utils} */
            // @ts-ignore
            const utils = window.utils;

            try {
                const resp = await fetch("{{ .ServerPathPrefix }}/api/ping");
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

    window.addEventListener("load", () => {
        window.dispatchEvent(new Event("focus"));
    });
})();
//{{ end }}
