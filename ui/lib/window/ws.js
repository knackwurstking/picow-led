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

        window.utils.setOnlineIndicatorState(false);

        // Reconnect here
        timeout = setTimeout(async () => {
            console.debug(`Try to reconnect to "${getURL()}"`);
            await ws.connect();
        }, timeoutDuration);
    };

    const onOpen = () => {
        console.debug(`WebSocket connected to "${getURL()}"`);
        window.utils.setOnlineIndicatorState(true);
    };

    /**
     * @param {MessageEvent<Blob>} ev
     * @returns {Promise<void>}
     */
    const onMessage = async (ev) => {
        /** @type {WSMessageData} */
        const data = JSON.parse(await ev.data.text());
        console.debug(`WebSocket message event:`, data);
        ws.events.dispatch(data.type, data.data);
    };

    /** @type {WS} */
    const ws = {
        events: new window.ui.Events(),

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
