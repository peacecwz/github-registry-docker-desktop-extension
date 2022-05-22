import React from 'react';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';

const PackageVersionSelector = () => {
    const { actions, state } = useStateMachine({ updateStore });

    return <div>
        {
            state?.github?.packageVersions?.length > 0
                ?
                <select onChange={e => actions.updateStore({
                    currentPackageVersion: e.target.value
                })}>
                    {state?.github?.packageVersions?.map(pkgVersion => <option value={pkgVersion.Name}>{pkgVersion.Name}</option>)}
                </select>
                : state?.github?.currentPackage && <h2>No package versions found</h2>
        }
        <br />
    </div>

}

export default PackageVersionSelector;