// src/components/common/PrivateRoute.js

import React from 'react';
import { Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';

function PrivateRoute({ children, roles }) {
    const auth = useSelector(state => state.auth);

    console.log('Auth state:', auth); 

    if (auth.loading) {
        // Покажите индикатор загрузки
        return <div>Loading...</div>;
    }

    if (!auth.isAuthenticated) {
        // User is not authenticated
        console.log("User is not authenticated", auth.user)
        return <Navigate to="/login" />;
    }

    if (roles && !roles.includes(auth.user.role_name)) {
        // User doesn't have the required role
        console.log("User doesn't have the required role", auth.user)
        return <Navigate to="/" />;
    }


    return children;
}

export default PrivateRoute;
