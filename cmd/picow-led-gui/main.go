package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func init() {
	log.SetFlags(0)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(newCustomTheme())

	w := a.NewWindow("PicoW LED")

	w.SetContent(newContent())
	w.Resize(fyne.NewSize(600, 800))
	w.ShowAndRun()
}
