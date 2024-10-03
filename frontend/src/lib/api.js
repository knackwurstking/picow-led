import { utils } from ".";

/**
 * @param {PicowStore} store
 * @param {string} path
 * @param {any} data
 * @returns {Promise<boolean>} ok
 */
export async function Delete(store, path, data) {
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

/**
 * @param {PicowStore} store
 * @param {string} path
 * @returns {Promise<any>} data - Returns `undefined` on error
 */
export async function Get(store, path) {
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

/**
 * @param {PicowStore} store
 * @param {string} path
 * @param {any} data
 * @returns {Promise<boolean>} ok
 */
export async function Post(store, path, data) {
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

/**
 * @param {PicowStore} store
 * @param {string} path
 * @param {any} data
 * @returns {Promise<boolean>} ok
 */
export async function Put(store, path, data) {
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

/**
 * @param {PicowStore} store
 * @param {string} path
 * @returns {Promise<string>}
 */
export async function getURL(store, path) {
    const server = store.ui.get("server");
    const addr = !server.port ? server.host : `${server.host}:${server.port}`;
    return `${server.ssl ? "https:" : "http:"}//${addr}${path}`;
}

/**
 * @param {Response} resp
 * @returns {Promise<boolean>} ok
 */
export async function handleResponseError(resp) {
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
