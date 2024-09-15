package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customColor interface {
	Variant() fyne.ThemeVariant
	SetVariant(variant fyne.ThemeVariant)
	RGBA() (r, g, b, a uint32)
}

type customColorBase struct {
	variant fyne.ThemeVariant
}

func (c *customColorBase) Variant() fyne.ThemeVariant {
	return c.variant
}

func (c *customColorBase) SetVariant(variant fyne.ThemeVariant) {
	c.variant = variant
}

type customColorBackground struct {
	customColorBase
}

func (c *customColorBackground) RGBA() (r, g, b, a uint32) {
	if c.Variant() == theme.VariantDark {
		return color.RGBA{9, 9, 11, 255}.RGBA() // NOTE: hsla(240, 10%, 4%, 1)
	}

	return color.RGBA{244, 244, 246, 255}.RGBA() // NOTE: hsla(240, 10%, 96%, 1)
}

type customColorForeground struct {
	customColorBase
}

func (c *customColorForeground) RGBA() (r, g, b, a uint32) {
	if c.Variant() == theme.VariantDark {
		return color.RGBA{227, 227, 232, 255}.RGBA() // NOTE: hsla(240, 10%, 90%, 1)
	}

	return color.RGBA{23, 23, 28, 255}.RGBA() // NOTE: hsla(240, 10%, 10%, 1)
}
