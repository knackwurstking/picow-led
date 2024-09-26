import { UIDrawer, html } from "ui";

/**
 * @typedef {{
 *   element: UIDrawer;
 * }} Drawer
 */

/**
 * @returns {Drawer}
 */
export default function () {
    const el = new UIDrawer();

    // TODO: Continue here...
    el.innerHTML = html``;

    return {
        element: el,
    };
}
