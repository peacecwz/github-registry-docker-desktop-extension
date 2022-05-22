import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { CssBaseline } from '@mui/material';
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import { StateMachineProvider, createStore } from 'little-state-machine';
import { createDockerDesktopClient } from '@docker/extension-api-client';

function log(store) {
  console.log(store);
}

createStore({
  isLoading: true,
  client: createDockerDesktopClient(),
  github: {
    currentOrganization: null,
    currentPackage: null,
    currentPackageVersion: null,
    organizations: [],
    packages: [],
    packageVersions: [],
    token: localStorage.getItem('token')// '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' // TODO (peacecwz): Remove hard coded token
  }
}, {
  middleWares: [log]
});

ReactDOM.render(
  <React.StrictMode>
    <DockerMuiThemeProvider>
      <CssBaseline />
      <StateMachineProvider>
        <App />
      </StateMachineProvider>
    </DockerMuiThemeProvider>
  </React.StrictMode>,
  document.getElementById("root")
);
