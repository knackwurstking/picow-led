package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct{}

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
}

func (*customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*customTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
