document.addEventListener("DOMContentLoaded", () => {
    // Starting the WebSocket handler
    window.ws.connect();

    window.ws.events.addListener("open", () => {
        console.debug("ws: open...");
        window.utils.setOnlineIndicatorState(true);
    });

    window.ws.events.addListener("close", () => {
        console.debug("ws: close...");
        window.utils.setOnlineIndicatorState(false);
    });

    window.ws.events.addListener("error", () => {
        console.debug("ws: error...");
    });

    window.ws.events.addListener("message", async (data) => {
        console.debug("ws: message...", data);

        switch (data.type) {
            case "colors":
                window.store.set("colors", data.data);
                break;
        }
    });
});
