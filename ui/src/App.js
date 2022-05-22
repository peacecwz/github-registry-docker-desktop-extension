import React, { useState } from "react";
import Button from "@mui/material/Button";
import { useStateMachine } from 'little-state-machine';
import updateStore from './store/change-store';
import PackageSelector from "./components/package-selector";
import PackageVersionSelector from "./components/package-version-selector";
import GithubDataTable from "./components/table";
import Layout from "./components/layout";
import DockerClient from './utils/docker-client'

const App = () => {
  const { actions, state } = useStateMachine({ updateStore });

  const onClose = async () => {
    await getOrganizations();
  }

  const getOrganizations = async () => {
    const result = await DockerClient.extension.vm.service.post("/organizations", {
      token: state?.github?.token
    });

    actions.updateStore({
      organizations: (result && result.Organizations) || []
    });
  };

  const pullPackage = async () => {
    const url = `ghcr.io/${state?.github?.currentOrganization}/${state?.github?.currentPackage}:${state?.github?.currentPackageVersion}`;

    try {
      var result = await DockerClient.docker.cli.exec("login", [
        'ghcr.io',
        '--username',
        state?.github?.currentOrganization,
        '--password',
        state?.github?.token
      ]);
      console.log('login', result.lines())
    } catch (e) {
      console.log('login-error', e)
    }

    var result = await DockerClient.docker.cli.exec("pull", [
      url
    ]);
    console.log('pull-package-version', result.lines())
  }

  const deletePackage = () => {
  }

  React.useEffect(() => {
    if (!window) return;

    getOrganizations();
  }, [])

  return (
    <Layout>
      <PackageSelector />
      <PackageVersionSelector />

      <GithubDataTable />

      <Button variant="contained" onClick={pullPackage}>
        Pull Package
      </Button>

      <Button variant="contained" onClick={deletePackage}>
        Delete Package
      </Button>
    </Layout>
  );
}

export default App;
