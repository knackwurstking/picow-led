import { UIAlert } from "ui";

/**
 * @param {import("ui").UIAlert_Options} options
 */
export function throwAlert(options) {
    /**
     * @type {import("ui").UIAlerts}
     */
    const alerts = document.querySelector(`ui-alerts`);
    if (!alerts) return;

    const alert = new UIAlert(options);
    alert.style.cursor = "pointer";

    const remove = alerts.ui.add(alert);
    alert.onclick = async () => remove();
}
