package components

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
)

var (
	// assetVersionCache stores the computed asset version to avoid recomputation
	assetVersionCache string
	// assetVersionComputed tracks if we've already computed the version
	assetVersionComputed bool
)

// AssetURL generates a versioned URL for an asset
func AssetURL(serverPathPrefix, assetPath string) templ.SafeURL {
	version := assetVersion()
	baseURL := serverPathPrefix + assetPath

	if version != "" {
		return templ.SafeURL(fmt.Sprintf("%s?v=%s", baseURL, version))
	}

	return templ.URL(baseURL)
}

// assetVersion returns a version string for cache-busting assets
func assetVersion() string {
	if assetVersionComputed {
		return assetVersionCache
	}

	// Generate version based on startup time
	// This ensures version changes on each deployment/restart
	startupTime := time.Now().Unix()
	assetVersionCache = fmt.Sprintf("%d", startupTime)
	assetVersionComputed = true
	return assetVersionCache
}
