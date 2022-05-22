import React, { useState } from "react";
import Button from "@mui/material/Button";
import { LoginPopup } from "./components/login-popup";
import { useStateMachine } from 'little-state-machine';
import updateStore from './store/change-store';
import PackageSelector from "./components/package-selector";
import PackageVersionSelector from "./components/package-version-selector";
import GithubDataTable from "./components/table";

function App() {
  const { actions, state } = useStateMachine({ updateStore });

  const [open, setOpen] = useState(state?.github?.token ? false : true);

  const onClose = async () => {
    console.log('onClose called')
    await getOrganizations();
  }

  const getOrganizations = async () => {
    const result = await state?.client?.extension.vm.service.post("/organizations", {
      token: state?.github?.token
    });

    actions.updateStore({
      organizations: result.Organizations || []
    });
  };

  const pullPackage = async () => {
    const url = `ghcr.io/${state?.github?.currentOrganization}/${state?.github?.currentPackage}:${state?.github?.currentPackageVersion}`;

    try {
      var result = await state?.client?.docker.cli.exec("login", [
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

    var result = await state?.client?.docker.cli.exec("pull", [
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
    <div className="App">
      <LoginPopup onClose={onClose} open={open} setOpen={setOpen} />
      <PackageSelector />
      <PackageVersionSelector />


      <GithubDataTable />

      <Button variant="contained" onClick={pullPackage}>
        Pull Package
      </Button>

      <Button variant="contained" onClick={deletePackage}>
        Delete Package
      </Button>
    </div>
  );
}

export default App;
