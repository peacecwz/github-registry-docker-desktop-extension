import React, { useState } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    TextField,
    DialogActions,
    Button
} from '@mui/material';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';

export const LoginPopup = ({ open, setOpen, onClose }) => {
    const { actions, state } = useStateMachine({ updateStore });
    const [token, setToken] = useState(state?.github?.token);

    const handleLogin = () => {
        localStorage.setItem('token', token);
        actions.updateStore({
            token
        });
        onClose && onClose();
        setOpen(false);
    }

    const handleClose = () => {
        setOpen(false);
    };

    return <div>
        <Dialog open={open} onClose={handleClose}>
            <DialogTitle>Login to the Github</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    Open Github Settings {'>'} Developer Settings {'>'} Personal access token {'>'} Generate New Token {'>'} Select {'write:packages '} {'delete:packages'} {' read:org'} and generate token
                </DialogContentText>
                <TextField
                    autoFocus
                    margin="dense"
                    id="personalAccessToken"
                    label="Github Personal Access Token"
                    type="password"
                    fullWidth
                    variant="standard"
                    onChange={(e) => setToken(e.target.value)}
                    value={token}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Cancel</Button>
                <Button onClick={handleLogin}>Login</Button>
            </DialogActions>
        </Dialog>
    </div>
}