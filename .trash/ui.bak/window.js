import store from "./window/store.js";
import api from "./window/api.js";
import ws from "./window/ws.js";
import utils from "./window/utils.js";

// @ts-expect-error
window.store = store();
// @ts-expect-error
window.api = api();
// @ts-expect-error
window.ws = ws();
// @ts-expect-error
window.utils = utils();
//!{{ end }}
