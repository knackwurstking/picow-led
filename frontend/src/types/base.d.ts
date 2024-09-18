/// <reference types="vite-plugin-pwa/client" />

interface Server {
  ssl: boolean;
  host: string;
  port: string;
}

interface Device {
  server: {
    name: string;
    addr: string;
    isOffline?: bool;
  };
  pins?: number[] | null;
  color?: number[] | null;
}
