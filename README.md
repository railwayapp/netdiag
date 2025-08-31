# Railway NetDiag

[**Download for macOS**](https://github.com/railwayapp/netdiag/releases/latest/download/RailwayNetDiag_macOS.zip) | [**Download for Windows**](https://github.com/railwayapp/netdiag/releases/latest/download/RailwayNetDiag_Win.zip)

macOS & Windows application for diagnosing client networking issues to Railway.

![railway-netdiag screenshot](./docs/screenshot.png)

## Usage

1. Download the latest version ([macOS](https://github.com/railwayapp/netdiag/releases/latest/download/RailwayNetDiag_macOS.zip) | [Windows](https://github.com/railwayapp/netdiag/releases/latest/download/RailwayNetDiag_Win.zip))
2. Open the application
3. Click "Run Diagnostics"
4. Wait for the diagnostics to complete
5. Copy the results to your clipboard or save them to a file
6. Share the results with Railway support for further assistance

## ⚠️ About Security Warnings

When you launch the application for the first time, you may encounter security
warnings from your OS. **This happens because the app is not code-signed yet**.
You can wait for a new release of the code-signed app (by mid-Sept '25), or
you can use it now by bypassing the warnings:

- On macOS:
  1. Open the app once
  2. Go to System Preferences -> Privacy & Security
  3. Scroll down to the warning at the bottom under "Allow applications from"
  4. Click "Open Anyway" for Railway NetDiag.app

- On Windows:
  1. Click "More info" on the "Windows protected your PC" popup
  2. Click "Run Anyway"

## Development

This is a [Wails](https://wails.io) + React application. To get started:

1. Clone the repository
2. Install [Wails](https://wails.io/docs/gettingstarted/installation) on your system
3. Run `wails dev` to start the application in development mode

### Structure

```
core/     - Go backend code
frontend/ - React frontend code
main.go   - Main entry point for the application
```

### Building

Use `wails build` to create a production build of the application:

```
wails build -clean -platform=windows/amd64     # Windows
wails build -clean -platform=darwin/universal  # macOS
```

This should only be necessary when testing in local dev. For production
releases, create a new release on GitHub and the build+release will be
performed automatically.

## License

[MIT Copyright (c) 2025 Railway Corporation](./LICENSE)
