package main

import (
	"embed"
	"encoding/json"
	"railway-network-debug/core"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed wails.json
var wailsJSON string

// WailsConfig represents the structure needed from wails.json
type WailsConfig struct {
	Info struct {
		ProductVersion string `json:"productVersion"`
	} `json:"info"`
}

func main() {
	// Parse version from wails.json
	var config WailsConfig
	version := "unknown"
	if err := json.Unmarshal([]byte(wailsJSON), &config); err == nil {
		if config.Info.ProductVersion != "" {
			version = config.Info.ProductVersion
		}
	}

	// Create an instance of the app structure
	app := core.NewApp()
	app.SetVersion(version)

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "Railway NetDiag",
		Width:         600,
		Height:        600,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
			&core.DiagnosticUpdate{},
		},
		EnumBind: []interface{}{
			core.AllDiagnosticUpdateTypes,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
