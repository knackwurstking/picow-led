package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	a.Settings().SetTheme(&customTheme{})

	w := a.NewWindow("PicoW LED")

	w.SetContent(NewContent())

	// TODO: Need to handle android size, text not centered with this custom theme
	w.Resize(fyne.NewSize(600, 800))

	w.ShowAndRun()
}
