import React from 'react';
import Header from './header';
import { useStateMachine } from 'little-state-machine';
import updateStore from '../store/change-store';

const Layout = ({ children }) => {
    const { state } = useStateMachine({ updateStore });
    return <div className={'App'}>
        <Header />
        {state?.github?.token ? <>{children}</> : <h3>You are not logged in</h3>}
    </div>
}

export default Layout;