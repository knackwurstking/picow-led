import { UIAlert, UIAlerts, type UIAlert_Options } from "ui";

export function throwAlert(options: UIAlert_Options) {
    const alerts = document.querySelector(`ui-alerts`) as UIAlerts;
    if (!alerts) return;

    const alert = new UIAlert(options);
    alert.style.cursor = "pointer";

    const remove = alerts.ui.add(alert);
    alert.onclick = async () => remove();
}
