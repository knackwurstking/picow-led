import type { PicowStore } from "../types";
import * as utils from "./utils";

export async function Delete(
    store: PicowStore,
    path: string,
    data: any
): Promise<boolean> {
    const url = await getURL(store, path);

    try {
        const resp = await fetch(url, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        return await handleResponseError(resp);
    } catch (err) {
        const message = `${url}: ${err}`;
        console.error(message);
        utils.throwAlert({ message, variant: "error" });
        return false;
    }
}

export async function Get(
    store: PicowStore,
    path: string
): Promise<any | undefined> {
    const url = await getURL(store, path);

    try {
        const resp = await fetch(url, { method: "GET" });

        const ok = await handleResponseError(resp);
        if (!ok) return undefined;
        return await resp.json();
    } catch (err) {
        const message = `${url}: ${err}`;
        console.error(message);
        utils.throwAlert({ message, variant: "error" });
        return undefined;
    }
}

export async function Post(
    store: PicowStore,
    path: string,
    data: any
): Promise<boolean> {
    const url = await getURL(store, path);

    try {
        const resp = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        return await handleResponseError(resp);
    } catch (err) {
        const message = `${url}: ${err}`;
        console.error(message);
        utils.throwAlert({ message, variant: "error" });
        return false;
    }
}

export async function Put(
    store: PicowStore,
    path: string,
    data: any
): Promise<boolean> {
    const url = await getURL(store, path);

    try {
        const resp = await fetch(url, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        return await handleResponseError(resp);
    } catch (err) {
        const message = `${url}: ${err}`;
        console.error(message);
        utils.throwAlert({ message, variant: "error" });
        return false;
    }
}

export async function getURL(store: PicowStore, path: string): Promise<string> {
    const server = store.ui.get("server");
    const addr = !server.port ? server.host : `${server.host}:${server.port}`;
    return `${server.ssl ? "https:" : "http:"}//${addr}${path}`;
}

export async function handleResponseError(resp: Response): Promise<boolean> {
    if (resp.ok) return true;

    resp.text().then((e) => {
        const message = `Server response to ${resp.url}: ${e}`;
        utils.throwAlert({ message, variant: "error" });
        console.error(message);
    });

    const message = `Fetch from "${resp.url}" with status code ${resp.status}`;
    console.error(message);
    utils.throwAlert({ message, variant: "error" });

    return false;
}
