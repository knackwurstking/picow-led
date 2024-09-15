package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	a.Settings().SetTheme(&customTheme{})

	w := a.NewWindow("PicoW LED")

	w.SetContent(NewContent())
	w.Resize(fyne.NewSize(600, 800))
	w.ShowAndRun()
}
