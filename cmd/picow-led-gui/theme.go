package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct{}

func newCustomTheme() fyne.Theme {
	return &customTheme{}
}

func (*customTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s) // TODO: Add monospaced fonts
	}

	if s.Bold {
		if s.Italic {
			return resourceRecMonoCasualBoldItalic1085Ttf
		}

		return resourceRecMonoCasualRegular1085Ttf
	}

	if s.Italic {
		return resourceRecMonoCasualItalic1085Ttf
	}

	return resourceRecMonoCasualRegular1085Ttf
}

func (*customTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
	//log.Printf("[DEBUG][*customTheme.Color] %s %d: %+v", n, v, color)
	//return color
}

func (*customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
	//log.Printf("[DEBUG][*customTheme.Icon] %s: %+v", n, icon.Name())
	//return icon
}

func (*customTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
	//log.Printf("[DEBUG][*customTheme.Size] %s: %+v", n, size)
	//return size
}
