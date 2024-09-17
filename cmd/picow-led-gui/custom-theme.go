package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	background customColor
	foreground customColor
	shadow     customColor
}

func newCustomTheme() fyne.Theme {
	return &customTheme{
		background: &customColorBackground{},
		foreground: &customColorForeground{},
		shadow:     &customColorShadow{},
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
	log.Printf("[DEBUG][*customTheme.Color] %s (%d)", n, v)

	switch n {
	case "background":
		c.background.SetVariant(v)
		return c.background
	case "foreground":
		c.foreground.SetVariant(v)
		return c.foreground
	case "shadow":
		c.shadow.SetVariant(v)
		return c.shadow
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
