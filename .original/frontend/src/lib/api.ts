import { addToast } from "./stores/toasts";

let pathPrefix = "";

export function setPathPrefix(prefix: string): void {
    pathPrefix = prefix;
}

export const PowerStateOFF = 0;
export const PowerStateON = 1;

export const devices = {
    async GET(): Promise<Devices> {
        const url = `${pathPrefix}/api/devices`;

        try {
            const r = await fetch(url);

            if (!r.ok) {
                throw new Error(r.statusText);
            }

            return await r.json();
        } catch (err) {
            addToast("error", `Fetch failed: ${url}: ${err}`);
            throw err;
        }
    },

    addr: {
        async GET(addr: string): Promise<Device> {
            const url = `${pathPrefix}/api/devices/${addr}`;

            try {
                const r = await fetch(url);

                if (!r.ok) {
                    throw new Error(r.statusText);
                }

                return await r.json();
            } catch (err) {
                addToast("error", `Fetch failed: ${url}: ${err}`);
                throw err;
            }
        },

        color: {
            async POST(addr: string, color: number[]): Promise<void> {
                const url = `${pathPrefix}/api/devices/${addr}/color`;

                try {
                    const r = await fetch(url, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify(color),
                    });

                    if (!r.ok) {
                        throw new Error(r.statusText);
                    }
                } catch (err) {
                    addToast("error", `Fetch failed: ${url}: ${err}`);
                    throw err;
                }
            },
        },

        power: {
            async GET(addr: string): Promise<PowerState> {
                const url = `${pathPrefix}/api/devices/${addr}/power`;

                try {
                    const r = await fetch(url);

                    if (!r.ok) {
                        throw new Error(r.statusText);
                    }

                    return await r.json();
                } catch (err) {
                    addToast("error", `Fetch failed: ${url}: ${err}`);
                    throw err;
                }
            },

            async POST(addr: string, state: PowerState) {
                const url = `${pathPrefix}/api/devices/${addr}/power?state=${state}`;

                try {
                    const r = await fetch(url, {
                        method: "POST",
                    });

                    if (!r.ok) {
                        throw new Error(r.statusText);
                    }
                } catch (err) {
                    addToast("error", `Fetch failed: ${url}: ${err}`);
                    throw err;
                }
            },
        },
    },
};

export const colors = {
    async GET(): Promise<Colors> {
        const url = `${pathPrefix}/api/colors`;

        try {
            const r = await fetch(url);

            if (!r.ok) {
                throw new Error(r.statusText);
            }

            return await r.json();
        } catch (err) {
            addToast("error", `Fetch failed: ${url}: ${err}`);
            throw err;
        }
    },

    async POST(colors: Colors): Promise<void> {
        const url = `${pathPrefix}/api/colors`;

        try {
            const r = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(colors),
            });

            if (!r.ok) {
                throw new Error(r.statusText);
            }
        } catch (err) {
            addToast("error", `Fetch failed: ${url}: ${err}`);
            throw err;
        }
    },

    async PUT(colors: Colors): Promise<void> {
        const url = `${pathPrefix}/api/colors`;

        try {
            const r = await fetch(url, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(colors),
            });

            if (!r.ok) {
                throw new Error(r.statusText);
            }
        } catch (err) {
            addToast("error", `Fetch failed: ${url}: ${err}`);
            throw err;
        }
    },

    id: {
        async POST(id: number, color: Color): Promise<void> {
            const url = `${pathPrefix}/api/colors/${id}`;

            try {
                const r = await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(color),
                });

                if (!r.ok) {
                    throw new Error(r.statusText);
                }
            } catch (err) {
                addToast("error", `Fetch failed: ${url}: ${err}`);
                throw err;
            }
        },
    },
};
