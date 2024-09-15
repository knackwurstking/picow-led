package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	background *colorBackground
	foreground *colorForeground
}

func newCustomTheme() fyne.Theme {
	return &customTheme{
		background: &colorBackground{},
		foreground: &colorForeground{},
	}
}

func (c *customTheme) Font(s fyne.TextStyle) fyne.Resource {
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

func (c *customTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	//log.Printf("[DEBUG][*customTheme.Color] %s (%d)", n, v)

	switch n {
	case "background":
		c.background.variant = v
		return c.background
	case "foreground":
		c.foreground.variant = v
		return c.foreground
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (c *customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
	//log.Printf("[DEBUG][*customTheme.Icon] %s: %+v", n, icon.Name())
	//return icon
}

func (c *customTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
	//log.Printf("[DEBUG][*customTheme.Size] %s: %+v", n, size)
	//return size
}

type colorBackground struct {
	variant fyne.ThemeVariant
}

func (c *colorBackground) RGBA() (r, g, b, a uint32) {
	if c.variant == theme.VariantDark {
		return color.RGBA{9, 9, 11, 255}.RGBA() // NOTE: hsla(240, 10%, 4%, 1)
	}

	return color.RGBA{244, 244, 246, 255}.RGBA() // NOTE: hsla(240, 10%, 96%, 1)
}

type colorForeground struct {
	variant fyne.ThemeVariant
}

func (c *colorForeground) RGBA() (r, g, b, a uint32) {
	if c.variant == theme.VariantDark {
		return color.RGBA{227, 227, 232, 255}.RGBA() // NOTE: hsla(240, 10%, 90%, 1)
	}

	return color.RGBA{23, 23, 28, 255}.RGBA() // NOTE: hsla(240, 10%, 10%, 1)
}
