const CACHE_VERSION = process.env.VERSION;
const CURRENT_CACHE = "picow-led-" + CACHE_VERSION;

const cacheFiles = [
    process.env.SERVER_PATH_PREFIX + "/",
    process.env.SERVER_PATH_PREFIX + "/settings",
    process.env.SERVER_PATH_PREFIX + "/css/style.css",
    process.env.SERVER_PATH_PREFIX + "/css/ui-v4.1.0.css",
    process.env.SERVER_PATH_PREFIX + "/js/ui-v4.1.0.min.umd.cjs",
    process.env.SERVER_PATH_PREFIX + "/js/devices-address.js",
    process.env.SERVER_PATH_PREFIX + "/js/devices.js",
    process.env.SERVER_PATH_PREFIX + "/js/layout.js",
    process.env.SERVER_PATH_PREFIX + "/js/main.js",
    process.env.SERVER_PATH_PREFIX + "/js/settings.js",
];

const blackList = ["/api/.*"];

self.addEventListener("activate", (evt) => {
    // @ts-expect-error
    evt.waitUntil(
        caches.keys().then((cacheNames) => {
            return Promise.all(
                cacheNames.map((cacheName) => {
                    if (cacheName !== CURRENT_CACHE) {
                        return caches.delete(cacheName);
                    }
                }),
            );
        }),
    );
});

self.addEventListener("install", (evt) => {
    // @ts-expect-error
    evt.waitUntil(
        (() => {
            caches
                .keys()
                .then((keyList) => {
                    return Promise.all(
                        keyList.map(function (key) {
                            return caches.delete(key);
                        }),
                    );
                })
                .then(() => {
                    caches.open(CURRENT_CACHE).then((cache) => {
                        return cache.addAll(cacheFiles);
                    });
                });
        })(),
    );
});

const fromNetwork = (
    /** @type {Request} */ request,
    /** @type {number} */ timeout,
) => {
    return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(reject, timeout);

        fetch(request).then((response) => {
            clearTimeout(timeoutId);
            update(request, response.clone());
            resolve(response);
        }, reject);
    });
};

const fromCache = async (/** @type {Request} */ request) => {
    const cache = await caches.open(CURRENT_CACHE);
    const matching = await cache.match(request);
    return matching || cache.match("/offline/");
};

const update = async (
    /** @type {Request} */ request,
    /** @type {Response} */ response,
) => {
    if (isBlackListed(request.url)) {
        return new Promise(() =>
            console.debug(`Nope, no caching for: ${request.url}`),
        );
    }

    const cache = await caches.open(CURRENT_CACHE);
    cache.put(request, response);
};

const isBlackListed = (/** @type {string} */ url) => {
    return blackList.find((path) => new RegExp(".*" + path + "$").test(url));
};

self.addEventListener(
    "fetch",
    (/** @type {Event & { request: Request }} */ evt) => {
        if (!isBlackListed(evt.request.url)) {
            // @ts-expect-error
            evt.respondWith(
                fromNetwork(evt.request, 1e4).catch(() =>
                    fromCache(evt.request),
                ),
            );
        }
    },
);
