/**
 * @typedef Api
 * @type {{
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
 *  setOnlineIndicator: (state: boolean) => void;
 *  powerButtonClickHandler: (ev: Event & { currentTarget: HTMLButtonElement }) => Promise<void>;
 *  registerServiceWorker: () => void;
 * }}
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
