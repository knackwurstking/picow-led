//{{ define "script-devices" }}
window.addEventListener("load", () => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    const items = utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"].innerText = "Devices";
});
//{{ end }}
