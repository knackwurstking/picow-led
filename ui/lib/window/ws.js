/**
 * @returns {WS}
 */
export function create() {
    let timeout = null;
    /** @type {number} */
    const timeoutDuration = 1000;

    function getURL() {
        return process.env.SERVER_PATH_PREFIX + `/ws`;
    }

    const onClose = () => {
        if (timeout !== null) {
            clearTimeout(timeout);
            timeout = null;
        }

        // Reconnect here
        timeout = setTimeout(async () => {
            console.debug(`Try to reconnect to "${getURL()}"`);
            await ws.connect();
        }, timeoutDuration);
    };

    const onOpen = () => {
        console.debug(`WebSocket connected to "${getURL()}"`);
    };

    /**
     * @param {MessageEvent<WSMessage>} ev
     * @returns {void}
     */
    const onMessage = (ev) => {
        console.debug(`WebSocket message event:`, ev);
        // TODO: ...
    };

    /** @type {WS} */
    const ws = {
        /** @type {WebSocket | null} */
        socket: null,

        isOpen() {
            if (!this.socket) return false;
            return this.socket.readyState === this.socket.OPEN;
        },

        async connect() {
            if (this.socket) this.close();

            const wsAddr = getURL(); // origin + path
            console.debug(`Try to connect WebSocket to ${wsAddr}`);

            this.socket = new WebSocket(wsAddr);

            // Reconnect handler
            this.socket.addEventListener("close", onClose);
            this.socket.addEventListener("open", onOpen);
            this.socket.addEventListener("message", onMessage);
        },

        close() {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }

            if (this.socket) {
                this.socket.removeEventListener("close", onClose);
                if (this.isOpen()) this.socket.close();
                this.socket = null;
            }
        },
    };

    return ws;
}
