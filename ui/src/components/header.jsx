import React, { useState } from 'react';
import {
    AppBar,
    Box,
    Toolbar,
    Typography,
    Button
} from '@mui/material';
import LoginPopup from './login-popup';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';

const Header = () => {
    const { state, actions } = useStateMachine({ updateStore });
    const [isLogged, setIsLogged] = useState(state?.github?.token ? true : false);
    const [open, setOpen] = useState(false);

    const onLogin = () => {
        setOpen(true);
    }

    const getIsLogged = () => {
        return isLogged;
    }

    const onLogout = () => {
        actions.updateStore({
            token: null
        });

        localStorage.clear();
        setIsLogged(false);
    }

    return <Box sx={{ flexGrow: 1 }}>
        <LoginPopup onClose={() => setIsLogged(true)} open={open} setOpen={setOpen} />
        <AppBar position="static">
            <Toolbar>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    Github Registry Manager
                </Typography>
                {getIsLogged()
                    ? <Button color="inherit" onClick={onLogout}>Logout</Button>
                    : <Button color="inherit" onClick={onLogin}>Login</Button>}
            </Toolbar>
        </AppBar>
    </Box>
}

export default Header;