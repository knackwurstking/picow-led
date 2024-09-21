import { UIAlert } from "ui";

/**
 * @param {object} options
 * @param {string} options.message
 */
export function throwAlert({ message }) {
    /**
     * @type {import("ui").UIAlerts}
     */
    const alerts = document.querySelector(`ui-alerts`);
    if (!alerts) return;

    const alert = new UIAlert({ message });
    const remove = alerts.ui.add(alert);
    alert.onclick = async () => remove();
}
