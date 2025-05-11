window.addEventListener("pageshow", () => {
    // Starting the WebSocket handler
    window.ws.connect();
});
