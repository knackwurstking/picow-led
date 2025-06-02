//go:build !frontend
// +build !frontend

package main

import (
	"embed"
)

func frontend() (f embed.FS, ok bool) {
	return embed.FS{}, false
}
