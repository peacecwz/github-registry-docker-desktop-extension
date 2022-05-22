import React from 'react';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';
import DockerClient from '../utils/docker-client';

const PackageSelector = () => {
    const { actions, state } = useStateMachine({ updateStore });

    const selectPackage = async (packageName) => {
        actions.updateStore({
            currentPackage: packageName
        })

        const result = await DockerClient.extension.vm.service.post(`/package-versions`, {
            packageName,
            token: state?.github?.token,
            organizationId: state?.github?.currentOrganization,
        });

        actions.updateStore({
            packageVersions: result.PackageVersions || []
        })
    }

    return <div>
        {
            state?.github?.packages?.length > 0
                ?
                <select onChange={async (e) => await selectPackage(e.target.value)}>
                    {state?.github?.packages?.map(pkg => <option value={pkg.Name}>{pkg.Name}</option>)}
                </select>
                : state?.github?.currentOrganization && <h2>No packages found</h2>
        }
        <br />
    </div>
}

export default PackageSelector;