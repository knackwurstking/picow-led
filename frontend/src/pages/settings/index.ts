import * as globals from "../../globals";

export async function onMount() {
    console.debug("[settings] Mount settings template/page");

    const routerTarget = document.querySelector(`.router-target`)!;
    const usingSSL = routerTarget.querySelector<HTMLInputElement>(`#settings-usingSSL`)!;
    const server = globals.store.get("server")!;

    let timeout: NodeJS.Timeout | null = null;

    const resetTimeout = () => {
        if (timeout != null) {
            clearTimeout(timeout);
            timeout = null;
        }
    };

    const update = (server: { ssl?: boolean; host?: string; port?: string }) => {
        timeout = setTimeout(() => {
            timeout = null;
            globals.store.update("server", (data) => {
                return {
                    ...data,
                    ...server,
                };
            });
        }, 250);
    };

    {
        usingSSL.checked = server.ssl;
        usingSSL.oninput = () => {
            resetTimeout();
            update({ ssl: usingSSL.checked });
        };
    }

    {
        const serverHost = routerTarget.querySelector<HTMLInputElement>(`#settings-serverHost`)!;
        serverHost.value = server.host;
        serverHost.oninput = () => {
            resetTimeout();
            update({ host: serverHost.value });
        };
    }

    {
        const serverPort = routerTarget.querySelector<HTMLInputElement>(`#settings-serverPort`)!;
        serverPort.value = server.port;
        serverPort.oninput = () => {
            resetTimeout();
            update({ port: serverPort.value });
        };
    }

    const devicesButton = document.querySelector<HTMLElement>(`.ui-app-bar button#goBack`)!;
    devicesButton.style.display = "block";
    devicesButton.onclick = () => {
        location.hash = "";
    };
}

export async function onDestroy() {
    console.debug("[settings] Destroy devices template/page");

    const devicesButton = document.querySelector<HTMLElement>(`.ui-app-bar button#goBack`)!;
    devicesButton.style.display = "none";
}
