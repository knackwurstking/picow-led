//!{{ define "script-page-settings" }}
(() => {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-ignore
    const w = window;

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = w.utils.setupAppBarItems(
            "online-indicator",
            "title",
            "back-button",
        );

        items["back-button"].onclick = (ev) => {
            ev.preventDefault();
            location.pathname = `{{ .ServerPathPrefix }}/`;
        };

        items["title"].innerText = "Settings";
    }

    window.addEventListener("load", () => {
        setupAppBar();
    });
})();
//!{{ end }}
