//{{ define "script-devices-addr" }}
window.addEventListener("load", async () => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    // Setup AppBar Items

    const items = utils.setupAppBarItems(
        "back-button",
        "online-indicator",
        "title",
    );

    items["back-button"].onclick = (ev) => {
        ev.preventDefault();
        location.pathname = `{{ .ServerPathPrefix }}/`;
    };

    items["title"].innerText = decodeURIComponent(
        location.pathname.split("/").reverse()[0],
    );
});
//{{ end }}
