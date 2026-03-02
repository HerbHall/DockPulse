import React from "react";
import ReactDOM from "react-dom/client";
import CssBaseline from "@mui/material/CssBaseline";
import { DockerMuiThemeProvider } from "@docker/docker-mui-theme";
import { ErrorBoundary } from "./components/ErrorBoundary";
import App from "./App";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <DockerMuiThemeProvider>
      <CssBaseline />
      <ErrorBoundary>
        <App />
      </ErrorBoundary>
    </DockerMuiThemeProvider>
  </React.StrictMode>,
);
