import { UIStore } from "../lib/store";
import { Api } from "../lib/api";
import * as utils from "../lib/window/utils";
import { WS } from "../lib/ws";

window.store = new UIStore();
window.api = new Api();
window.utils = utils.create();
window.ws = new WS();
