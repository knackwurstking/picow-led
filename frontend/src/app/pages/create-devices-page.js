/**
 * @typedef {{
 *  element: UIStackLayoutPage;
 *  get name(): string;
 * }} SettingsPage
 */

import { UIStackLayoutPage } from "ui";

/**
 * @param {object} options
 * @param {PicowStore} options.store
 * @returns {Promise<SettingsPage>}
 */
export default async function ({ store }) {
    const pageName = "devices";
    const page = new UIStackLayoutPage(pageName);

    // TODO: Create the "Devices" page here

    return {
        element: page,

        get name() {
            return pageName;
        },
    };
}
