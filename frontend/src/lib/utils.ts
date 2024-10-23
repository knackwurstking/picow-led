import { UIAlert, UIAlerts, UIAlertVariant } from "ui";

export function throwAlert(options: {
    variant: UIAlertVariant;
    message: string;
}) {
    const alerts = document.querySelector<UIAlerts>(`ui-alerts`)!;

    const alert = new UIAlert();

    alert.variant = options.variant;
    alert.message = options.message;

    alert.style.cursor = "pointer";

    const cleanUp = alerts.addAlert(alert);
    alert.onclick = async () => cleanUp();
}
