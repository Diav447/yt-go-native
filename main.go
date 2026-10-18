package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.yt.gonative")
	w := a.NewWindow("YT Go Native - No WebView")
	w.Resize(fyne.NewSize(400,700))
	entry := widget.NewEntry()
	entry.SetPlaceHolder("ketik kurumi terus ENTER di keyboard")
	label := widget.NewLabel("hasil akan muncul di sini")
	entry.OnSubmitted = func(s string){
		label.SetText("Searching: " + s + " ... (API fix next)")
	}
	w.SetContent(container.NewVBox(entry, label))
	w.ShowAndRun()
}
