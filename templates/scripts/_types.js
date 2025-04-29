/**
 * @typedef UI
 * @type {import("ui")}
 *
 * @typedef Api
 * @type {{
 *  devices: () => Promise<Device[]>;
 *  setDevicesColor: (color: MicroColor | undefined | null, ...devices: Device[]) => Promise<Device[]>
 * }}
 *
 * @typedef WS
 * @type {{
 *  addr: () => string;
 *  isOpen: () => boolean;
 *  connect: () => void;
 *  close: () => void;
 * }}
 *
 * @typedef Utils
 * @type {{
 *  onClickPowerButton: (ev: Event & { currentTarget: HTMLButtonElement }) => Promise<void>;
 *  updateDeviceListItem: (item: HTMLElement, device: Device) => void;
 *  setupAppBarItems: (...itemNames: AppBarItemName[]) => AppBarItems;
 *  setOnlineIndicatorState: (state: boolean) => void;
 *  registerServiceWorker: () => void;
 * }}
 *
 * @typedef AppBarItemName
 * @type {"back-button" | "online-indicator" | "title" | "settings-button"}
 *
 * @typedef AppBarItems
 * @type {{ [key: string]: HTMLElement }}
 *
 * @typedef Device
 * @type {{
 *  server: Server
 *  online: boolean
 *  error: string
 *  color: MicroColor
 *  pins: MicroPins
 * }}
 *
 * @typedef Server
 * @type {{
 *  addr: string
 *  name: string
 * }}
 *
 * @typedef {number[]} MicroColor
 * @typedef {number[]} MicroPins
 */
