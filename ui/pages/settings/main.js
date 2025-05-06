(() => {
    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = window.utils.setupAppBarItems(
            "online-indicator",
            "title",
        );
        items["title"].innerText = "Settings";
    }

    window.addEventListener("pageshow", () => {
        setupAppBar();
    });
})();
