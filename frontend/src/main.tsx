import React from "react";
import { createRoot } from "react-dom/client";
import "./style.css";
import App from "./App";

const container = document.getElementById("root");
if (!container) throw new Error("no root element!");
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);

// Open all links externally
document.body.addEventListener("click", (e: MouseEvent) => {
  if (!e.target) return;
  if (!(e.target instanceof HTMLAnchorElement)) return;
  if (e.target.nodeName !== "A") return;
  if (!e.target.href) return;

  const url = e.target.href;
  if (url.startsWith("http://wails.localhost:")) return;
  if (url.startsWith("file://")) return;
  if (url.startsWith("http://#")) return;

  e.preventDefault();
  // @ts-expect-error wails runtime is injected in window, but not in types
  window.runtime.BrowserOpenURL(url);
});
