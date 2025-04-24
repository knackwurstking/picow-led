/**
 * Returns the newly created alert element
 */
export function add(variant: "info" | "warning" | "error", message: string): HTMLElement {
    const alertsContainer = document.querySelector(`.alerts-container`)!;

    const alert = document.createElement("div");
    alertsContainer.appendChild(alert);

    alert.style.padding = "var(--ui-spacing)";
    alert.style.margin = "0.25rem";

    switch (variant) {
        case "warning":
            alert.style.color = "var(--ui-warning-text)";
            alert.style.backgroundColor = "var(--ui-warning)";
            break;
        case "error":
            alert.style.color = "var(--ui-error-text)";
            alert.style.backgroundColor = "var(--ui-error)";
            break;
        default:
            alert.style.color = "var(--ui-info-text)";
            alert.style.backgroundColor = "var(--ui-info)";
    }

    alert.className = `alert ${variant} ui-flex-grid-item`;
    alert.innerHTML = `<p><i>${message}</i></p>`;
    alert.onclick = () => {
        alert.parentElement!.removeChild(alert);
    };

    return alert;
}
