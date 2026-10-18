package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Video struct {
	Title   string `json:"title"`
	VideoID string `json:"videoId"`
}

func main() {
	a := app.NewWithID("com.yt.gonative")
	w := a.NewWindow("YT Go Native")
	w.Resize(fyne.NewSize(400, 700))

	var videos []Video
	list := widget.NewList(
		func() int { return len(videos) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(videos[i].Title)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		if id < len(videos) {
			// buka di app youtube native, bukan webview
			link := fmt.Sprintf("https://youtu.be/%s", videos[id].VideoID)
			u, _ := url.Parse(link)
			a.OpenURL(u)
		}
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("ketik kurumi lalu ENTER")
	status := widget.NewLabel("Ketik lalu tekan ENTER di keyboard HP")

	entry.OnSubmitted = func(q string) {
		status.SetText("Searching: " + q + "...")
		go func() {
			resp, err := http.Get("https://inv.nadeko.net/api/v1/search?q=" + url.QueryEscape(q))
			if err != nil {
				status.SetText("Error: " + err.Error())
				return
			}
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&videos)
			if len(videos) == 0 {
				status.SetText("Ga ketemu")
			} else {
				status.SetText(fmt.Sprintf("Ketemu %d video", len(videos)))
			}
			list.Refresh()
		}()
	}

	w.SetContent(container.NewBorder(entry, status, nil, nil, list))
	w.ShowAndRun()
}
