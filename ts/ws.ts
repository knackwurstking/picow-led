// TODO: WebSocket handler here

class WS {
    private socket: WebSocket | null = null;

    private timeout: number | null = null;
    private timeoutDuration: number = 1000;

    private closeHandler = async () => {
        if (this.timeout !== null) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        // Reconnect here
        this.timeout = setTimeout(async () => {
            this.connect();
        }, this.timeoutDuration);
    };

    public addr(): string {
        return ``; // TODO: ...
    }

    protected isOpen(): boolean {
        if (!this.socket) return false;
        return this.socket.readyState === this.socket.OPEN;
    }

    protected connect(): void {
        if (this.socket) this.close();

        const wsAddr = this.addr(); // origin + path
        console.debug(`Try to connect WebSocket to ${wsAddr}`);

        this.socket = new WebSocket(wsAddr);

        // Reconnect handler
        this.socket.addEventListener("close", this.closeHandler);
    }

    protected close() {
        if (!!this.timeout) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        if (this.socket) {
            this.socket?.removeEventListener("close", this.closeHandler);
            if (this.isOpen()) this.socket.close();
            this.socket = null;
        }
    }
}
