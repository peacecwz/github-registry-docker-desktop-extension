import React from 'react';
import { Autocomplete, TextField } from '@mui/material';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';

const OrganizationSelector = () => {
    const { actions, state } = useStateMachine({ updateStore });

    const selectOrg = async (orgId) => {
        actions.updateStore({
            currentOrganization: orgId
        });

        const result = await state?.client?.extension.vm.service.post(`/packages`, {
            token: state?.github?.token,
            organizationId: orgId
        });

        console.log('result', result);

        actions.updateStore({
            packages: result.Packages || []
        })
    }

    return <Autocomplete
        disablePortal
        onChange={async (_, v) => await selectOrg(v.label)}
        options={state?.github?.organizations?.map((org) => ({
            label: org.Name,
            id: org.Id
        }))}
        sx={{ width: 300 }}
        renderInput={(params) => <TextField {...params} />}
    />
}

export default OrganizationSelector;