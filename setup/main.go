package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed payload/Slate.exe
var payloadSlate []byte

func main() {
	app := NewApp(payloadSlate)
	err := wails.Run(&options.App{
		Title:            "Win-Slate Setup",
		Width:            760,
		Height:           680,
		MinWidth:         680,
		MinHeight:        560,
		BackgroundColour: &options.RGBA{R: 12, G: 13, B: 16, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind:      []interface{}{app},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
