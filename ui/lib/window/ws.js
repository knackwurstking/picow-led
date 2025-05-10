/**
 * @returns {import("../../types").WS}
 */
export function create() {
    /** @type {WebSocket | null} */
    let socket = null;
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
            console.debug(`Try to reconnect to "${getURL()}"`);
            connect();
        }, timeoutDuration);
    };

    const onOpen = () => {
        console.debug(`WebSocket connected to "${getURL()}"`);
    };

    /**
     * @param {MessageEvent<import("../../types").WSMessage>} ev
     * @returns {void}
     */
    const onMessage = (ev) => {
        // TODO: Continue here
        console.debug(`WebSocket message event:`, ev);
    };

    function getURL() {
        return process.env.SERVER_PATH_PREFIX + `/ws`;
    }

    function isOpen() {
        if (!socket) return false;
        return socket.readyState === socket.OPEN;
    }

    function connect() {
        if (socket) close();

        const wsAddr = getURL(); // origin + path
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

    /** @type {import("../../types").WS} */
    return {
        isOpen,
        connect,
        close,
    };
}
