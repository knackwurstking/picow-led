/**
 * @param {PicowStore} store
 * @param {string} path
 * @returns {Promise<string>}
 */
export async function url(store, path) {
    const server = store.ui.get("server");
    const addr = !server.port ? server.host : `${server.host}:${server.port}`;
    return `${server.ssl ? "https:" : "http:"}//${addr}${path}`;
}
