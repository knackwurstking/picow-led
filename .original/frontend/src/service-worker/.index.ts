/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />

import { build, files, version } from "$service-worker";

const CURRENT_CACHE = `picow-led-${version}`;

const cacheFiles = [
    ...build,
    ...files,

    "https://fonts.googleapis.com/css2?family=Recursive:slnt,wght,CASL,CRSV,MONO@-15..0,300..1000,0..1,0..1,0..1&display=swap",
];

// NOTE: This will be prefixed with ".*" and suffixed with "$"
//const blackList = ["/screenshots/.*", "/service-worker.js", "/manifest.json"];
const blackList = [];

self.addEventListener("activate", (evt) => {
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

const fromNetwork = (request: Request, timeout: number) => {
    return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(reject, timeout);

        fetch(request).then((response) => {
            clearTimeout(timeoutId);
            update(request, response.clone());
            resolve(response);
        }, reject);
    });
};

const fromCache = async (request: Request) => {
    const cache = await caches.open(CURRENT_CACHE);
    const matching = await cache.match(request);
    return matching || cache.match("/offline/");
};

const update = async (request: Request, response: Response) => {
    const cache = await caches.open(CURRENT_CACHE);
    cache.put(request, response);
};

const isBlackListed = (url: string) => {
    return blackList.find((path) => new RegExp(".*" + path + "$").test(url));
};

// @ts-expect-error
self.addEventListener("fetch", (evt: Event & { request: Request }) => {
    if (!isBlackListed(evt.request.url)) {
        console.debug(`fetch: caching: ${evt.request.url}`);

        // @ts-expect-error
        evt.respondWith(
            fromNetwork(evt.request, 1e4).catch(() => fromCache(evt.request)),
        );
    } else {
        console.debug(`fetch: no caching for: ${evt.request.url}`);
    }
});
