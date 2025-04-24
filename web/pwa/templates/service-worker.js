// This is the service worker code that lives at the root

// You have to supply a name for your cache, this will
// allow us to remove an old one to avoid hitting disk
// space limits and displaying old resources
var cacheName = "picow-led-{{ .Version }}";

// Assets to cache
var assetsToCache = [
    {{ $prefix := .ServerPathPrefix }}
    {{ range .AssetsToCache }}
        "{{ $prefix }}/{{ . }}",
    {{ end }}
];

self.addEventListener("install", function (event) {
    // waitUntil() ensures that the Service Worker will not
    // install until the code inside has successfully occurred
    event.waitUntil(
        // Create cache with the name supplied above and
        // return a promise for it
        caches
            .open(cacheName)
            .then(function (cache) {
                // Important to `return` the promise here to have `skipWaiting()`
                // fire after the cache has been updated.
                return cache.addAll(assetsToCache);
            })
            .then(function () {
                // `skipWaiting()` forces the waiting ServiceWorker to become the
                // active ServiceWorker, triggering the `onactivate` event.
                // Together with `Clients.claim()` this allows a worker to take effect
                // immediately in the client(s).
                return self.skipWaiting();
            }),
    );
});

// Activate event
// Be sure to call self.clients.claim()
self.addEventListener("activate", function () {
    // `claim()` sets this worker as the active worker for all clients that
    // match the workers scope and triggers an `oncontrollerchange` event for
    // the clients.
    return self.clients.claim();
});

self.addEventListener("fetch", function (event) {
    // Ignore non-get request like when accessing the admin panel
    if (event.request.method !== "GET") {
        return;
    }

    console.debug(`fetch: ${event.request.url}`, event.request);

    // Don't try to handle non-secure assets because fetch will fail
    //if (/http:/.test(event.request.url)) {
    //  return;
    //}

    // Here's where we cache all the things!
    event.respondWith(
        // Open the cache created when install
        caches.open(cacheName).then(async function (cache) {
            // Go to the network to ask for that resource
            try {
                const networkResponse = await fetch(event.request);
                // Add a copy of the response to the cache (updating the old version)
                // NOTE: Stuff like "chrome-extension" will maybe trigger an error here
                //       because "chrome-extension" is unsupported, i will ignore
                //       for now.
                cache.put(event.request, networkResponse.clone());
                return networkResponse;
            } catch {
                return await cache.match(event.request);
            }
        }),
    );
});
