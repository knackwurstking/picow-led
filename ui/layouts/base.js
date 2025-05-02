//!{{ define "script-layout-base" }}
(() => {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-ignore
    const w = window;

    /** @type {import("../types.d.ts").UIStore} */
    const store = new w.ui.Store("picow-led:");

    /** @type {number | null} */
    let timeout = null;
    window.addEventListener("focus", () => {
        if (timeout !== null) {
            return;
        }

        // NOTE: This focus event gets triggered twice after tab switching,
        //       this is just a workaround to prevent this
        timeout = setTimeout(async () => {
            try {
                const resp = await fetch("{{ .ServerPathPrefix }}/api/ping");
                const data = await resp.text();
                if (data === "pong") {
                    w.utils.setOnlineIndicatorState(true);
                    updateStorage();
                } else {
                    w.utils.setOnlineIndicatorState(false);
                }
            } catch (err) {
                console.error(err);
                w.utils.setOnlineIndicatorState(false);
            }

            timeout = null;
        }, 100);
    });

    window.addEventListener("load", () => {
        // Trigger the focus event once
        window.dispatchEvent(new Event("focus"));
    });

    function updateStorage() {
        w.api.devices().then((devices) => {
            store.set("devices", devices);
        });
    }
})();
//!{{ end }}
