package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Application menu — mirrors Electron About / Help / Support links
	// plus a thin Brain menu for Claude Code / Codex CLI login (no API keys).
	appMenu := menu.NewMenu()
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		if app.ctx != nil {
			runtime.Quit(app.ctx)
		}
	})

	brainMenu := appMenu.AddSubmenu("Brain")
	brainMenu.AddText("Connect Claude account…", nil, func(_ *menu.CallbackData) {
		app.MenuConnectClaude()
	})
	brainMenu.AddText("Test Claude brain…", nil, func(_ *menu.CallbackData) {
		app.MenuTestClaude()
	})
	brainMenu.AddSeparator()
	brainMenu.AddText("Connect Codex account…", nil, func(_ *menu.CallbackData) {
		app.MenuConnectCodex()
	})
	brainMenu.AddSeparator()
	brainMenu.AddText("Refresh brain status…", nil, func(_ *menu.CallbackData) {
		app.MenuRefreshBrain()
	})
	brainMenu.AddText("Claude Code docs", nil, func(_ *menu.CallbackData) {
		app.OpenExternal("https://claude.com/claude-code")
	})

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("About Slate", nil, func(_ *menu.CallbackData) {
		app.EmitAbout()
	})
	helpMenu.AddText("Slate Help", keys.CmdOrCtrl("/"), func(_ *menu.CallbackData) {
		app.EmitHelp()
	})
	helpMenu.AddSeparator()
	helpMenu.AddText("Support Slate on Ko-fi ♥", nil, func(_ *menu.CallbackData) {
		app.OpenExternal("https://ko-fi.com/samwasserman")
	})
	helpMenu.AddText("wassermanproductions.com", nil, func(_ *menu.CallbackData) {
		app.OpenExternal("https://wassermanproductions.com")
	})
	helpMenu.AddText("wasserman.ai", nil, func(_ *menu.CallbackData) {
		app.OpenExternal("https://wasserman.ai")
	})
	helpMenu.AddSeparator()
	helpMenu.AddText("Original project (GitHub)", nil, func(_ *menu.CallbackData) {
		app.OpenExternal("https://github.com/wassermanproductions/slate")
	})

	err := wails.Run(&options.App{
		Title:     "Win-Slate",
		Width:     1520,
		Height:    940,
		MinWidth:  1100,
		MinHeight: 680,
		Menu:      appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 12, G: 13, B: 16, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
