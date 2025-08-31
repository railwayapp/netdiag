package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	version string
}

func NewApp() *App {
	return &App{
		version: "unknown",
	}
}

// SetVersion sets the application version
func (a *App) SetVersion(version string) {
	a.version = version
}

// GetAppVersion returns the application version
func (a *App) GetAppVersion() string {
	return a.version
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// DiagnosticUpdateType represents the type of diagnostic update
type DiagnosticUpdateType string

// Diagnostic update types
const (
	DIAG_TYPE_START         DiagnosticUpdateType = "start"
	DIAG_TYPE_STEP_START    DiagnosticUpdateType = "step_start"
	DIAG_TYPE_STEP_PROGRESS DiagnosticUpdateType = "step_progress"
	DIAG_TYPE_DONE          DiagnosticUpdateType = "done"
	DIAG_TYPE_ERROR         DiagnosticUpdateType = "error"
)

// DiagnosticUpdate represents a streaming update from diagnostics
type DiagnosticUpdate struct {
	Type    DiagnosticUpdateType `json:"type"`
	Message string               `json:"message"`
	Data    string               `json:"data"`
}

// AllDiagnosticUpdateTypes contains all possible diagnostic update type values for enum binding
var AllDiagnosticUpdateTypes = []struct {
	Value  DiagnosticUpdateType
	TSName string
}{
	{DIAG_TYPE_START, "START"},
	{DIAG_TYPE_STEP_START, "STEP_START"},
	{DIAG_TYPE_STEP_PROGRESS, "STEP_PROGRESS"},
	{DIAG_TYPE_DONE, "DONE"},
	{DIAG_TYPE_ERROR, "ERROR"},
}

// emitDiagnosticsOutputEvent emits a diagnostic output event to wails runtime
func (a *App) emitDiagnosticsOutputEvent(update DiagnosticUpdate) {
	if update.Type == DIAG_TYPE_STEP_START {
		sep := strings.Repeat("-", 79)
		update.Data = "\n" + sep + "\n" + update.Data + "\n" + sep
	}
	if update.Type == DIAG_TYPE_START || update.Type == DIAG_TYPE_STEP_START || update.Type == DIAG_TYPE_STEP_PROGRESS {
		update.Data = update.Data + "\n"
	}
	runtime.EventsEmit(a.ctx, "diag-output", update)
}

// RunDiagnosticsStream executes diagnostic commands and streams output via
// wails events
func (a *App) RunDiagnosticsStream() {
	now := time.Now().Format("Monday, Jan 2 2006 15:04:05 MST")
	genOnMsg := fmt.Sprintf("Generated : %s", now)
	startMsg := "Railway Network Diagnostics\n"
	startMsg += genOnMsg + "\n"
	startMsg += "Endpoint  : " + RAILWAY_ROUTING_INFO_ENDPOINT + "\n"
	startMsg += "\n"

	go func() {

		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_START,
			Message: "Starting diagnostics...",
			Data:    strings.TrimSuffix(startMsg, "\n"),
		})

		time.Sleep(100 * time.Millisecond)

		// 1. IPInfo
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Fetching client IP info...",
			Data:    "Client IP Info",
		})
		ipData, ipErr := a.runIPInfoDiagnostic()
		if ipErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "IP info failed",
				Data:    fmt.Sprintf("Error: %v\n", ipErr),
			})
		} else if ipErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "Client IP information retrieved",
				Data:    ipData,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// 2. HTTP HEAD request
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Making HTTP request...",
			Data:    "HTTP HEAD request",
		})
		httpData, httpErr := a.runHttpHeadRequestDiagnostic()
		if httpErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "HTTP request failed",
				Data:    fmt.Sprintf("Error: %v\n", httpErr),
			})
		} else if httpErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "HTTP request completed",
				Data:    httpData,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// 3. Dig with system DNS
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Running DNS lookup (system DNS)...",
			Data:    "DNS lookup (using system DNS)",
		})
		digSystemData, digSystemErr := a.runDigSystemDNSDiagnostic()
		if digSystemErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "DNS lookup (system) failed",
				Data:    fmt.Sprintf("Error: %v\n", digSystemErr),
			})
			if digSystemData != "" {
				a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
					Type:    DIAG_TYPE_STEP_PROGRESS,
					Message: "Partial output",
					Data:    digSystemData,
				})
			}
		} else if digSystemErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "DNS lookup (system) completed",
				Data:    digSystemData,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// 4. Dig with Cloudflare DNS
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Running DNS lookup (Cloudflare)...",
			Data:    "DNS lookup (using Cloudflare)",
		})
		digCloudflareData, digCloudflareErr := a.runDigCloudflareDNSDiagnostic()
		if digCloudflareErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "DNS lookup (Cloudflare) failed",
				Data:    fmt.Sprintf("Error: %v\n", digCloudflareErr),
			})
			if digCloudflareData != "" {
				a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
					Type:    DIAG_TYPE_STEP_PROGRESS,
					Message: "Partial output",
					Data:    digCloudflareData,
				})
			}
		} else if digCloudflareErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "DNS lookup (Cloudflare) completed",
				Data:    digCloudflareData,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// 5. Traceroute
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Running traceroute (this may take a while)...",
			Data:    "Traceroute",
		})
		tracerouteErr := a.runTracerouteDiagnostic(func(line string) {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "Traceroute running (this may take a while)...",
				Data:    line,
			})
		})
		if tracerouteErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "Traceroute failed",
				Data:    fmt.Sprintf("Error: %v\n", tracerouteErr),
			})
		} else if tracerouteErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "Traceroute completed",
				Data:    "",
			})
		}
		time.Sleep(100 * time.Millisecond)

		// 6. Ping
		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_STEP_START,
			Message: "Running ping test (n=10)...",
			Data:    "Ping (n=10)",
		})
		pingErr := a.runPingDiagnostic(func(line string) {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "Ping running (this may take a awhile)...",
				Data:    line,
			})
		})
		if pingErr != nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_ERROR,
				Message: "Ping test failed",
				Data:    fmt.Sprintf("Error: %v\n", pingErr),
			})
		} else if pingErr == nil {
			a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
				Type:    DIAG_TYPE_STEP_PROGRESS,
				Message: "Ping completed",
				Data:    "",
			})
		}

		a.emitDiagnosticsOutputEvent(DiagnosticUpdate{
			Type:    DIAG_TYPE_DONE,
			Message: "Diagnostics complete!",
			Data:    "\nCompleted",
		})
	}()
}

// GetDiagnosticUpdateSchema returns an empty DiagnosticUpdate to expose the struct for TypeScript generation
func (a *App) GetDiagnosticUpdateSchema() DiagnosticUpdate {
	return DiagnosticUpdate{}
}

// SaveToFile saves the diagnostic output to a file
func (a *App) SaveToFile(content string) error {
	filename, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: fmt.Sprintf("railway-netdiag-%s.txt", time.Now().Format("2006-01-02-150405")),
		Title:           "Save Diagnostics Report",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return fmt.Errorf("error opening save dialog: %v", err)
	}

	if filename == "" {
		// User cancelled
		return nil
	}

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}
