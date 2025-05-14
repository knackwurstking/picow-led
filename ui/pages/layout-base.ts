window.addEventListener("pageshow", () => {
    // Starting the WebSocket handler
    window.ws.connect();

    let timeout: NodeJS.Timeout | null = null;
    window.addEventListener("focus", () => {
        if (timeout === null) {
            timeout = setTimeout(async () => {
                try {
                    if (!window.ws.isOpen()) await window.ws.connect();
                } finally {
                    timeout = null;
                }
            });
        }
    });
});
