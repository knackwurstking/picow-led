import * as store from "../lib/window/store";
import { Api } from "../lib/api";
import * as utils from "../lib/window/utils";
import * as ws from "../lib/window/ws";

window.store = store.create();
window.api = new Api();
window.utils = utils.create();
window.ws = ws.create();
