/**
 * @returns {import("../types.d.ts").UIStore}
 */
export function create() {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-expect-error
    const w = window;

    /** @type {import("../types.d.ts").UIStore} */
    const store = new w.ui.Store("picow-led:");

    store.set("devices", [], true);

    return store;
}
