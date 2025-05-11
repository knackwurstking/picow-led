import * as store from "../lib/window/store";
import * as api from "../lib/window/api";
import * as utils from "../lib/window/utils";
import * as ws from "../lib/window/ws";

window.store = store.create();
window.api = api.create();
window.utils = utils.create();
window.ws = ws.create();
