setOnlineIndicator("online");

function setOnlineIndicator(state: "online" | "offline") {
    const onlineIndicator =
        document.querySelector<HTMLElement>(`.online-indicator`)!;
    onlineIndicator.setAttribute(`data-state`, state);
}
