import React from 'react';
import { Autocomplete, TextField } from '@mui/material';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';
import DockerClient from '../utils/docker-client';

const OrganizationSelector = () => {
    const { actions, state } = useStateMachine({ updateStore });

    const selectOrg = async (orgId) => {
        actions.updateStore({
            currentOrganization: orgId
        });

        const result = await DockerClient.extension.vm.service.post(`/packages`, {
            token: state?.github?.token,
            organizationId: orgId
        });

        actions.updateStore({
            packages: result.packages || []
        })
    }

    return <Autocomplete
        disablePortal
        onChange={async (_, v) => await selectOrg(v.label)}
        options={state?.github?.organizations?.map((org) => ({
            label: org.name,
            id: org.id
        }))}
        sx={{ width: 300 }}
        renderInput={(params) => <TextField {...params} />}
    />
}

export default OrganizationSelector;