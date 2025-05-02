//!{{ define "script-window" }}
/** @type {import("../types.d.ts").PageWindow} */
// @ts-expect-error
const w = window;

const store = new w.ui.Store("picow-led:");

w.store.set("devices", [], true);

// @ts-ignore
window.store = store;
