//{{ define "script-base-layout" }}
(() => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    window.addEventListener("focus", async () => {
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
    });
})();
//{{ end }}
