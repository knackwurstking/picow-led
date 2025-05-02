(() => {
    // @ts-ignore
    window.store = require("./window/store").create();
    // @ts-ignore
    window.api = require("./window/api").create();
    // @ts-ignore
    window.ws = require("./window/ws").create();
    // @ts-ignore
    window.utils = require("./window/utils").create();
})();
