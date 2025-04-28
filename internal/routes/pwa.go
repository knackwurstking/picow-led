package routes

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

const manifest = `
{
    "version": "{{ .Version }}",
    "manifest_version": 3,
    "name": "PicoW LED",
    "short_name": "picow-led",
    "id": "",
    "icons": [
        {
            "src": "{{ .ServerPathPrefix }}/static/icons/pwa-64x64.png",
            "sizes": "64x64",
            "type": "image/png"
        },
        {
            "src": "{{ .ServerPathPrefix }}/static/icons/pwa-192x192.png",
            "sizes": "192x192",
            "type": "image/png"
        },
        {
            "src": "{{ .ServerPathPrefix }}/static/icons/pwa-512x512.png",
            "sizes": "512x512",
            "type": "image/png"
        },
        {
            "src": "{{ .ServerPathPrefix }}/static/icons/maskable-icon-512x512.png",
            "sizes": "512x512",
            "type": "image/png",
            "purpose": "maskable"
        }
    ],
    "screenshots": [
        {
            "src": "{{ .ServerPathPrefix }}/static/screenshots/626x338.png",
            "sizes": "626x338",
            "type": "image/png",
            "form_factor": "wide",
            "label": "App Preview"
        },
        {
            "src": "{{ .ServerPathPrefix }}/static/screenshots/328x626.png",
            "sizes": "328x626",
            "type": "image/png",
            "form_factor": "narrow",
            "label": "App Preview"
        }
    ],
    "theme_color": "#09090b",
    "background_color": "#09090b",
    "display": "standalone",
    "scope": ".",
    "start_url": "./",
    "public_path": "{{ .ServerPathPrefix }}"
}
`

// TODO: Move cacheFiles to routes.Options
const serviceWorker = `
const CACHE_VERSION = "{{ .Version }}";
const CURRENT_CACHE = "picow-led-" + "${CACHE_VERSION}";

const cacheFiles = [
    "{{ .ServerPathPrefix }}",
    "{{ .ServerPathPrefix }}/settings",
    "{{ .ServerPathPrefix }}/css/style.css",
    "{{ .ServerPathPrefix }}/css/ui-v4.1.0.css",
    "{{ .ServerPathPrefix }}/js/api.js",
    "{{ .ServerPathPrefix }}/js/utils.js",
    "{{ .ServerPathPrefix }}/js/ws.js",
];

self.addEventListener("activate", (evt) => {
    console.debug("activate:", evt);

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
    console.debug("install:", evt);

    evt.waitUntil(
        caches.open(CURRENT_CACHE).then((cache) => {
            return cache.addAll(cacheFiles);
        }),
    );
});

const fromNetwork = (request, timeout) => {
    return new Promise((fulfill, reject) => {
        const timeoutId = setTimeout(reject, timeout);
        fetch(request).then((response) => {
            clearTimeout(timeoutId);
            fulfill(response);
            update(request);
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

const update = (request) => {
    return caches
        .open(CURRENT_CACHE)
        .then((cache) =>
            fetch(request).then((response) => cache.put(request, response)),
        );
};

self.addEventListener("fetch", (evt) => {
    console.debug("fetch:", evt);

    evt.respondWith(
        fromNetwork(evt.request, 1e4).catch(() => fromCache(evt.request)),
    );

    evt.waitUntil(update(evt.request));
});
`

func pwa(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"/manifest.json", func(c echo.Context) error {
		t, err := template.New("manifest.json").Parse(manifest)
		if err != nil {
			return err
		}
		c.Response().Header().Add("Content-Type", "application/json")
		return t.Execute(c.Response().Writer, data)
	})

	e.GET(data.ServerPathPrefix+"/js/service-worker.js", func(c echo.Context) error {
		t, err := template.New("service-worker.js").Parse(serviceWorker)
		if err != nil {
			return err
		}
		c.Response().Header().Add("Content-Type", "application/javascript")
		c.Response().Header().Add("Service-Worker-Allowed", "/")
		return t.Execute(c.Response().Writer, data)
	})
}
