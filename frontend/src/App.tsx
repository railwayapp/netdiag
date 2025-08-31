import { useEffect, useRef, useState } from "react";
import logo from "../assets/images/logo.png";
import {
  GetAppVersion,
  RunDiagnosticsStream,
  SaveToFile,
} from "../wailsjs/go/core/App";
import { core } from "../wailsjs/go/models";
import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import { Button } from "./components/button";

// Event name for emitted diagnostic output
const EVT_DIAG_OUTPUT = "diag-output";

const App = () => {
  const [diagnosticOutput, setDiagnosticOutput] = useState("");
  const [isRunning, setIsRunning] = useState(false);
  const [statusMessage, setStatusMessage] = useState("");
  const [copied, setCopied] = useState(false);
  const [appVersion, setAppVersion] = useState("unknown");
  const outputRef = useRef<HTMLDivElement>(null);

  // Fetch app version on mount
  useEffect(() => {
    const getAppVersion = async () => {
      const version = await GetAppVersion();
      setAppVersion(version);
    };
    getAppVersion();
  }, []);

  // Auto-scroll to bottom when new content is added
  useEffect(() => {
    if (!diagnosticOutput) return;
    if (diagnosticOutput === "") return;
    if (!outputRef.current) return;
    outputRef.current.scrollTop = outputRef.current.scrollHeight;
  }, [diagnosticOutput]);

  // Ctrl/Cmd+A to select only output text
  useEffect(() => {
    const handleSelectAll = (e: KeyboardEvent) => {
      if (!(e.metaKey || e.ctrlKey)) return;
      if (e.key !== "a") return;
      if (!outputRef.current) return;
      if (!diagnosticOutput) return;
      if (diagnosticOutput === "") return;

      e.preventDefault();
      const range = document.createRange();
      range.selectNodeContents(outputRef.current);
      const selection = window.getSelection();
      selection?.removeAllRanges();
      selection?.addRange(range);
    };
    document.addEventListener("keydown", handleSelectAll);
    return () => document.removeEventListener("keydown", handleSelectAll);
  }, [diagnosticOutput]);

  // Process diagnostic output events
  useEffect(() => {
    EventsOn(EVT_DIAG_OUTPUT, (update: core.DiagnosticUpdate) => {
      setStatusMessage(update.message);
      setDiagnosticOutput((prev) => prev + update.data);
      if (update.type === core.DiagnosticUpdateType.DONE) {
        setIsRunning(false);
        setStatusMessage("");
      }
      if (update.type === core.DiagnosticUpdateType.ERROR) {
        setIsRunning(false);
        setStatusMessage(`Error: ${update.message}`);
      }
    });
    return () => {
      EventsOff(EVT_DIAG_OUTPUT);
    };
  }, []);

  const runDiagnostics = async () => {
    setIsRunning(true);
    setDiagnosticOutput("");
    setStatusMessage("Starting diagnostics...");

    try {
      await RunDiagnosticsStream();
    } catch (error) {
      setDiagnosticOutput(`Error starting diagnostics: ${error}`);
      setIsRunning(false);
      setStatusMessage("");
    }
  };

  const copyToClipboard = async () => {
    if (!diagnosticOutput) return;
    try {
      await navigator.clipboard.writeText(diagnosticOutput);
      setCopied(true);
      setTimeout(() => setCopied(false), 1000);
    } catch (error) {
      alert(`Failed to copy to clipboard: ${error}`);
      setCopied(false);
    }
  };

  const saveToFile = async () => {
    if (!diagnosticOutput) {
      alert("No diagnostics output to save");
      return;
    }
    try {
      await SaveToFile(diagnosticOutput);
    } catch (error) {
      alert(`Failed to save file: ${error}`);
    }
  };

  return (
    <div
      id="App"
      style={{
        padding: "10px 30px",
        margin: "0 auto",
      }}
    >
      <header
        className="noselect"
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <h1 style={{ margin: 0 }}>Railway Network Diagnostics</h1>
        <img
          src={logo}
          alt="Railway NetDiag Logo"
          style={{ width: "150px", height: "150px", pointerEvents: "none" }}
        />
      </header>

      <div className="noselect">
        <p style={{ color: "var(--gray-700)" }}>
          A network diagnostic tool that runs a series of tests to Railway's
          diagnostic endpoints and provides an output that can be shared with
          Railway support.
        </p>
      </div>

      {/* Status bar and controls */}
      <div
        className="noselect"
        style={{
          height: "60px",
          borderRadius: "12px 12px 0 0",
          background:
            "linear-gradient(135deg, var(--pink-300) 0%, var(--pink-400) 33%, var(--pink-500) 66%, var(--pink-600) 100%)",
          boxShadow:
            "0 2px 8px rgba(0, 0, 0, 0.3), inset 0 1px 0 rgba(255, 255, 255, 0.1)",
          border: "1px solid rgba(255, 255, 255, 0.1)",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "0 20px",
          position: "relative",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "10px",
            minWidth: "200px",
          }}
        >
          {isRunning ? (
            <>
              <div
                className="spinner"
                style={{
                  width: "14px",
                  height: "14px",
                  border: "2px solid var(--pink-300)",
                  borderTop: "3px solid var(--pink-800)",
                  borderRadius: "50%",
                  animation: "spin 1s linear infinite",
                }}
              ></div>
              <span style={{ color: "var(--pink-700)" }}>{statusMessage}</span>
            </>
          ) : (
            <Button onClick={runDiagnostics} disabled={isRunning}>
              Run Diagnostics
            </Button>
          )}
        </div>

        <div style={{ flex: 1 }}></div>

        {isRunning ||
          (diagnosticOutput && (
            <div style={{ display: "flex", gap: "10px", alignItems: "center" }}>
              <Button
                onClick={copyToClipboard}
                disabled={!diagnosticOutput || isRunning}
              >
                {copied ? "👍 Copied" : "Copy to Clipboard"}
              </Button>
              <Button
                onClick={saveToFile}
                disabled={!diagnosticOutput || isRunning}
              >
                Save to File
              </Button>
            </div>
          ))}
      </div>

      {/* Output */}
      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          backgroundColor: "var(--gray-100)",
          borderRadius: "5px",
          overflow: "hidden",
        }}
      >
        <div
          ref={outputRef}
          style={{
            color: "#f0f0f0",
            padding: "25px",
            overflowY: "auto",
            fontFamily: "monospace",
            fontSize: "10px",
            whiteSpace: "pre-wrap",
            wordBreak: "break-word",
            minHeight: "20vh",
            maxHeight: "20vh",
          }}
        >
          {diagnosticOutput || "Click 'Run Diagnostics' to start..."}
        </div>
      </div>

      {/* Footer */}
      <footer className="noselect" style={{ marginTop: "40px" }}>
        <p
          style={{
            color: "var(--gray-500)",
            fontSize: "12px",
            textAlign: "center",
          }}
        >
          &copy; Railway Corporation 2025 &bull; Version {appVersion} &bull;{" "}
          <a
            style={{ color: "var(--gray-500)" }}
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/railwayapp/netdiag"
          >
            Source Code
          </a>
        </p>
      </footer>
    </div>
  );
};

export default App;
