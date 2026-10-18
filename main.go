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
  "github.com/kkdai/youtube/v2"
)

type SearchResult struct {
  Items []struct{
    Title string `json:"title"`
    URL string `json:"url"`
    Thumbnail string `json:"thumbnail"`
  } `json:"items"`
}

func main(){
  a := app.New()
  w := a.NewWindow("YT Go Native")
  w.Resize(fyne.NewSize(380, 700))

  searchEntry := widget.NewEntry()
  searchEntry.SetPlaceHolder("Cari video... MrBeast, windah, dll")
  
  list := widget.NewList(
    func() int { return 0 },
    func() fyne.CanvasObject { return widget.NewLabel("template") },
    func(i widget.ListItemID, o fyne.CanvasObject) {},
  )
  
  var videos []struct{Title, ID string}
  
  searchEntry.OnSubmitted = func(q string){
    // 1. Search pake Piped (proxy YT biar ga kena bot)
    resp, _ := http.Get("https://pipedapi.kavin.rocks/search?q="+url.QueryEscape(q))
    defer resp.Body.Close()
    var data map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&data)
    // parse hasil
    videos = nil
    if items, ok := data["items"].([]interface{}); ok {
      for _, it := range items {
        if m, ok := it.(map[string]interface{}); ok {
          if m["type"]=="stream" {
            videos = append(videos, struct{Title, ID string}{
              Title: fmt.Sprint(m["title"]),
              ID: fmt.Sprint(m["url"]),
            })
          }
        }
      }
    }
    list.Length = func() int { return len(videos) }
    list.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject){
      o.(*widget.Label).SetText(videos[i].Title)
    }
    list.Refresh()
  }

  list.OnSelected = func(id widget.ListItemID){
    // 2. Narik stream asli YT
    go func(){
      client := youtube.Client{}
      vid, _ := client.GetVideo(videos[id].ID)
      stream := vid.Formats.WithAudioChannels()[0]
      // 3. Play native - disini nanti pake fyne video player
      fmt.Println("PLAY URL:", stream.URL) // ini URL mp4 asli dari googlevideo.com
      // Untuk Fyne, tinggal: player := widget.NewVideoFromURL(url.Parse(stream.URL))
    }()
  }

  content := container.NewBorder(
    container.NewVBox(widget.NewLabelWithStyle("YT Go Native - No WebView", fyne.TextAlignCenter, fyne.TextStyle{Bold:true}), searchEntry),
    nil, nil, nil, list,
  )
  w.SetContent(content)
  w.ShowAndRun()
}
