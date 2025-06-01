export class UIStore extends window.ui.Store<{
    colors: Colors;
    currentDeviceColors: Record<string, Color>;
}> {
    constructor() {
        super("picow-led:");

        this.set("colors", [], true);

        this.set("currentDeviceColors", {}, true);

        // @ts-expect-error
        this.delete("color"); // TODO: Just to clean up, can be removed before release
        // @ts-expect-error
        this.delete("devices"); // TODO: Just to clean up, can be removed before release
    }

    public currentDeviceColor(addr: string): Color | null {
        return this.get("currentDeviceColors")?.[addr] || null;
    }
}
