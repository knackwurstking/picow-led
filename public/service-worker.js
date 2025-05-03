const CACHE_VERSION = "v1.0.0";
const CURRENT_CACHE = "picow-led-" + CACHE_VERSION;
const cacheFiles = ["/", "/settings", "/css/style.css", "/css/ui-v4.1.0.css", "/js/ui-v4.1.0.min.umd.cjs", "/js/devices-address.js", "/js/devices.js", "/js/layout.js", "/js/main.js", "/js/settings.js"];
const blackList = ["/api/.*"];
self.addEventListener("activate", (evt) => {
  evt.waitUntil(caches.keys().then((cacheNames) => {
    return Promise.all(cacheNames.map((cacheName) => {
      if (cacheName !== CURRENT_CACHE) {
        return caches.delete(cacheName);
      }
    }));
  }));
});
self.addEventListener("install", (evt) => {
  evt.waitUntil((() => {
    caches.keys().then((keyList) => {
      return Promise.all(keyList.map(function(key) {
        return caches.delete(key);
      }));
    }).then(() => {
      caches.open(CURRENT_CACHE).then((cache) => {
        return cache.addAll(cacheFiles);
      });
    });
  })());
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
const fromCache = async (request) => {
  const cache = await caches.open(CURRENT_CACHE);
  const matching = await cache.match(request);
  return matching || cache.match("/offline/");
};
const update = async (request, response) => {
  if (isBlackListed(request.url)) {
    return new Promise(() => console.debug(`Nope, no caching for: ${request.url}`));
  }
  const cache = await caches.open(CURRENT_CACHE);
  cache.put(request, response);
};
const isBlackListed = (url) => {
  return blackList.find((path) => new RegExp(".*" + path + "$").test(url));
};
self.addEventListener("fetch", (evt) => {
  if (!isBlackListed(evt.request.url)) {
    evt.respondWith(fromNetwork(evt.request, 1e4).catch(() => fromCache(evt.request)));
  }
});
