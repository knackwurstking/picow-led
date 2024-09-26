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

    el.innerHTML = html``;

    return {
        element: el,
    };
}
