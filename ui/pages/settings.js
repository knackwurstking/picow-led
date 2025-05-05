(() => {
    const utils = require("../lib/utils");

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = utils.setupAppBarItems("online-indicator", "title");
        items["title"].innerText = "Settings";
    }

    window.addEventListener("pageshow", () => {
        setupAppBar();
    });
})();
