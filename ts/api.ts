interface Device {}

class Api {
    private url(): string {
        return `{{ .ServerPathPrefix }}`;
    }

    // NOTE: No need for this right now, i'm doing some server side rendering for now
    //
    //public async devices(): Promise<Device[]> {
    //    const url = this.url() + "/api/devices";
    //
    //    // TODO: GET "/api/devices"
    //
    //    return [];
    //}
}

(window as any).api = new Api();
