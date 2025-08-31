package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	goruntime "runtime"
	"strings"
	"time"
)

const RAILWAY_ROUTING_INFO_ENDPOINT = "routing-info-production.up.railway.app"

// runIPInfoDiagnostic fetches IP information
func (a *App) runIPInfoDiagnostic() (string, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return "", fmt.Errorf("IP info failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading IP info response: %w", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		return string(body), nil
	}

	return prettyJSON.String(), nil
}

// runHttpHeadRequestDiagnostic performs HTTP HEAD request
func (a *App) runHttpHeadRequestDiagnostic() (string, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		// Don't follow redirects - we want to see the original response headers
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", "https://"+RAILWAY_ROUTING_INFO_ENDPOINT, nil)
	if err != nil {
		return "", fmt.Errorf("creating HTTP request: %w", err)
	}

	// Add a user agent
	req.Header.Set("User-Agent", "Railway-Network-Debug/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making HTTP HEAD request: %w", err)
	}
	defer resp.Body.Close()

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
	result.WriteString(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
	result.WriteString(fmt.Sprintf("Protocol: %s\n\n", resp.Proto))

	result.WriteString("Response Headers:\n")

	// Sort headers for consistent output
	headerNames := make([]string, 0, len(resp.Header))
	for name := range resp.Header {
		headerNames = append(headerNames, name)
	}

	for _, name := range headerNames {
		values := resp.Header[name]
		for _, value := range values {
			result.WriteString(fmt.Sprintf("  %s: %s\n", name, value))
		}
	}

	return result.String(), nil
}

// runTracerouteDiagnostic runs traceroute and streams output
func (a *App) runTracerouteDiagnostic(callback TracerouteStreamCallback) error {
	return runTracerouteStream(callback)
}

// runPingDiagnostic runs ping test and streams output
func (a *App) runPingDiagnostic(callback PingStreamCallback) error {
	return runPingStream(callback)
}

// runDigSystemDNSDiagnostic runs dig using system DNS
func (a *App) runDigSystemDNSDiagnostic() (string, error) {

	output, err := runDigCommand("", RAILWAY_ROUTING_INFO_ENDPOINT)
	if err != nil {
		return output, fmt.Errorf("dig (system DNS) failed: %w", err)
	}

	return output, nil
}

// runDigCloudflareDNSDiagnostic runs dig using Cloudflare DNS
func (a *App) runDigCloudflareDNSDiagnostic() (string, error) {

	output, err := runDigCommand("1.1.1.1", RAILWAY_ROUTING_INFO_ENDPOINT)
	if err != nil {
		return output, fmt.Errorf("dig (Cloudflare DNS) failed: %w", err)
	}

	return output, nil
}

// TracerouteStreamCallback is called for each line of traceroute output
type TracerouteStreamCallback func(line string)

// runTracerouteStream runs a traceroute command and streams output via callback
func runTracerouteStream(callback TracerouteStreamCallback) error {
	var cmd *exec.Cmd
	if goruntime.GOOS == "windows" {
		// Windows tracert with 3 second timeout
		cmd = exec.Command("tracert", "-w", "3000", RAILWAY_ROUTING_INFO_ENDPOINT)
	} else {
		// macOS and Linux:
		//   -I   - Use ICMP
		//   -n   - Do not resolve hostnames
		//   -q 2 - Number of queries per hop
		//   -w 3 - Wait time per probe in seconds
		cmd = exec.Command("traceroute", "-I", "-q", "2", "-w", "3", RAILWAY_ROUTING_INFO_ENDPOINT)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting traceroute: %w", err)
	}

	var hasOutput bool
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		hasOutput = true
		line := scanner.Text()
		callback(line)
	}

	if err := cmd.Wait(); err != nil {
		// Traceroute often exits with non-zero status, which is normal
		// Only return error if there was no output
		if !hasOutput {
			// Only show stderr if there was an actual failure
			stderrBytes, _ := io.ReadAll(stderr)
			if len(stderrBytes) > 0 {
				callback("STDERR: " + string(stderrBytes))
			}
			return fmt.Errorf("traceroute command failed: %w", err)
		}
	}

	return nil
}

// runDigCommand runs a dig command (or nslookup on Windows) and returns the output
func runDigCommand(dnsServer string, hostname string) (string, error) {
	var cmd *exec.Cmd

	if goruntime.GOOS == "windows" {
		// Windows: use nslookup
		if dnsServer != "" {
			// Use specific DNS server
			cmd = exec.Command("nslookup", hostname, dnsServer)
		} else {
			// Use system DNS
			cmd = exec.Command("nslookup", hostname)
		}
	} else {
		// macOS and Linux: use dig
		if dnsServer != "" {
			// Use specific DNS server (e.g., @1.1.1.1)
			cmd = exec.Command("dig", "@"+dnsServer, hostname)
		} else {
			// Use system DNS
			cmd = exec.Command("dig", hostname)
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if command is not installed
		if strings.Contains(err.Error(), "executable file not found") {
			if goruntime.GOOS == "windows" {
				return "nslookup command not available on this system\n", nil
			} else {
				return "dig command not available on this system\n", nil
			}
		}
		return string(output), fmt.Errorf("DNS lookup command failed: %w", err)
	}

	return string(output), nil
}

// PingStreamCallback is called for each line of ping output
type PingStreamCallback func(line string)

// runPingStream runs a ping command and streams output via callback
func runPingStream(callback PingStreamCallback) error {
	var cmd *exec.Cmd

	if goruntime.GOOS == "windows" {
		// Windows: ping -n 5
		cmd = exec.Command("ping", "-n", "10", RAILWAY_ROUTING_INFO_ENDPOINT)
	} else {
		// macOS and Linux: ping -c 5
		cmd = exec.Command("ping", "-c", "10", RAILWAY_ROUTING_INFO_ENDPOINT)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting ping: %w", err)
	}

	var hasOutput bool
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			hasOutput = true
			callback(line)
		}
	}

	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		callback("STDERR: " + string(stderrBytes))
	}

	if err := cmd.Wait(); err != nil {
		if !hasOutput {
			return fmt.Errorf("ping command failed: %w", err)
		}
	}

	return nil
}
