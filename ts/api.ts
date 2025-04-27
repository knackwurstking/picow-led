interface Device {}

class Api {
    private url(): string {
        return `{{ .ServerPathPrefix }}`;
    }

    public async devices(): Promise<Device[]> {
        const url = this.url() + "/api/devices";

        // TODO: GET "/api/devices"

        return [];
    }
}

(window as any).api = new Api();
