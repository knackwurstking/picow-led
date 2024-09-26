type PicowStore = import("ui").UIStore<PicowStore_Events>;

interface PicowStore_Events {
  devices: Device[];
  currentPage: PicowStackLayout_Pages;
  server: Server;
}
