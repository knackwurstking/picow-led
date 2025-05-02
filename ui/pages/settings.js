(() => {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-ignore
    const w = window;

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = w.utils.setupAppBarItems("online-indicator", "title");
        items["title"].innerText = "Settings";
    }

    window.addEventListener("pageshow", () => {
        setupAppBar();
    });
})();
