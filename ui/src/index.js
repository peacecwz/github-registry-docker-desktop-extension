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

const getDockerDesktopClient = () => {
  try {
    return createDockerDesktopClient();
  } catch {
    // TODO (peacecwz): Return mock Docker Desktop Client or connect with NodeJS Dev script on Websocket
    return {
      extension: {
        vm: {
          service: {
            get: async (url) => {
              console.log(`[DOCKER_DESKTOP_CLIENT GET]: Url: ${url}`)
              return Promise.resolve();
            },
            post: async (url, data) => {
              console.log(`[DOCKER_DESKTOP_CLIENT POST]: Url: ${url} Data:`, data)
              return Promise.resolve();
            }
          },
        },
      },
    };
  }
}

createStore({
  isLoading: false,
  client: getDockerDesktopClient(),
  github: {
    currentOrganization: null,
    currentPackage: null,
    currentPackageVersion: null,
    organizations: [],
    packages: [],
    packageVersions: [],
    token: localStorage.getItem('token')
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
