import { UIAlert } from "ui";

/**
 * @param {"error" | "info"} variant
 * @param {object} options
 * @param {string} options.message
 */
export function throwAlert(variant, { message }) {
    /**
     * @type {import("ui").UIAlerts}
     */
    const alerts = document.querySelector(`ui-alerts`);
    if (!alerts) return;

    const alert = new UIAlert({ message });
    alert.ui.variant = variant;
    alert.style.cursor = "pointer";

    const remove = alerts.ui.add(alert);
    alert.onclick = async () => remove();
}
