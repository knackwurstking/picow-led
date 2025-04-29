//{{ define "script-window-ws" }}
(() => {
    /** @type {WebSocket | null} */
    let socket = null;
    /** @type {number | null} */
    let timeout = null;
    /** @type {number} */
    const timeoutDuration = 1000;

    const onClose = function () {
        if (timeout !== null) {
            clearTimeout(timeout);
            timeout = null;
        }

        // Reconnect here
        timeout = setTimeout(() => {
            connect();
        }, timeoutDuration);
    };

    const onOpen = () => {
        // TODO: Request devices
        //const devices = window.api.devices();
    };

    const onMessage = () => {
        // TODO: ...
    };

    function addr() {
        return ``; // TODO: ...
    }

    function isOpen() {
        if (!socket) return false;
        return socket.readyState === socket.OPEN;
    }

    function connect() {
        if (socket) close();

        const wsAddr = addr(); // origin + path
        console.debug(`Try to connect WebSocket to ${wsAddr}`);

        socket = new WebSocket(wsAddr);

        // Reconnect handler
        socket.addEventListener("close", onClose);
        socket.addEventListener("open", onOpen);
        socket.addEventListener("message", onMessage);
    }

    function close() {
        if (timeout) {
            clearTimeout(timeout);
            timeout = null;
        }

        if (socket) {
            socket.removeEventListener("close", onClose);
            if (isOpen()) socket.close();
            socket = null;
        }
    }

    /** @type {WS} */
    const ws = {
        addr,
        isOpen,
        connect,
        close,
    };

    window.ws = ws;
})();
//{{ end }}
