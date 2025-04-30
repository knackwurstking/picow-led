//{{ define "script-settings" }}
window.addEventListener("load", () => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    const items = utils.setupAppBarItems(
        "online-indicator",
        "title",
        "back-button",
    );

    items["back-button"].onclick = (ev) => {
        ev.preventDefault();
        location.pathname = `{{ .ServerPathPrefix }}/`;
    };

    items["title"].innerText = "Settings";
});
//{{ end }}
