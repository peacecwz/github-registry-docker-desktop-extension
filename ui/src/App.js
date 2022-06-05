import React from "react";
import { useStateMachine } from "little-state-machine";
import updateStore from "./store/change-store";
import Layout from "./components/layout";
import DockerClient from "./utils/docker-client";
import PackagesTable from "./components/packages";

const App = () => {
  const { actions, state } = useStateMachine({ updateStore });

  const pullPackage = async () => {
    const url = `ghcr.io/${state?.github?.currentOrganization}/${state?.github?.currentPackage}:${state?.github?.currentPackageVersion}`;

    try {
      var result = await DockerClient.docker.cli.exec("login", [
        "ghcr.io",
        "--username",
        state?.github?.currentOrganization,
        "--password",
        state?.github?.token,
      ]);
      console.log("login", result.lines());
    } catch (e) {
      console.log("login-error", e);
    }

    var result = await DockerClient.docker.cli.exec("pull", [url]);
    console.log("pull-package-version", result.lines());
  };

  const deletePackage = () => {};

  return (
    <Layout>
      <PackagesTable />
    </Layout>
  );
};

export default App;
