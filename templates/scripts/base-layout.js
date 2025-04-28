//{{ define "script-base-layout" }}
(() => {
    const utils = window.utils;

    window.addEventListener("focus", async () => {
        try {
            const resp = await fetch("{{ .ServerPathPrefix }}/api/ping");
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
//{{ end }}
