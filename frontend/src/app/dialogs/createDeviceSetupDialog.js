import { Events, UIDialog } from "ui";

/**
 * @typedef DeviceSetupDialog_Events
 * @type {import("ui").UIDialog_Events & { "submit": Device, "delete": Device }}
 *
 * @typedef {{
 *  element: UIDialog;
 *  events: Events<DeviceSetupDialog_Events>;
 * }} DeviceSetupDialog
 */

/**
 * @returns {Promise<DeviceSetupDialog>}
 */
export default async function ({}) {
    // TODO: Pass parameter: name, addr, pins, allowDeletion
    const dialog = new UIDialog("Device Setup");

    /**
     * @type {Events<DeviceSetupDialog_Events>}
     */
    const events = new Events();

    // TODO: ...

    return {
        element: dialog,
        events,
    };
}
