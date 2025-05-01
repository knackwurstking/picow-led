const CACHE_VERSION = "{{ .Version }}";
const CURRENT_CACHE = "picow-led-" + CACHE_VERSION;

const cacheFiles = [
    "{{ .ServerPathPrefix }}/",
    "{{ .ServerPathPrefix }}/settings",
    "{{ .ServerPathPrefix }}/css/style.css",
    "{{ .ServerPathPrefix }}/css/ui-v4.1.0.css",
    "{{ .ServerPathPrefix }}/js/ui-v4.1.0.min.umd.cjs",
];

const blackList = ["/api/.*"];

self.addEventListener("activate", (evt) => {
    //console.debug("activate:", evt);

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
    //console.debug("install:", evt);

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

const fromNetwork = (request, timeout) => {
    return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(reject, timeout);

        fetch(request).then((response) => {
            clearTimeout(timeoutId);
            update(request, response.clone());
            resolve(response);
        }, reject);
    });
};

const fromCache = (request) => {
    return caches
        .open(CURRENT_CACHE)
        .then((cache) =>
            cache
                .match(request)
                .then((matching) => matching || cache.match("/offline/")),
        );
};

const update = (request, response) => {
    if (
        blackList.find((path) =>
            new RegExp(".*" + path + "$").test(request.url),
        )
    ) {
        return new Promise(() =>
            console.warn(`Nope, no caching for: ${request.url}`),
        );
    }

    return caches.open(CURRENT_CACHE).then((cache) => {
        cache.put(request, response);
    });
};

self.addEventListener("fetch", (evt) => {
    //console.debug("fetch:", evt);

    evt.respondWith(
        fromNetwork(evt.request, 1e4).catch(() => fromCache(evt.request)),
    );
});
