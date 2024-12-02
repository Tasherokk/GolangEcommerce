// src/components/common/Navbar.js

import React from 'react';
import { Link } from 'react-router-dom';
import { Menu } from 'antd';
import { useSelector, useDispatch } from 'react-redux';
import { logout } from '../../store/actions/authActions';

function Navbar() {
  const auth = useSelector((state) => state.auth);
  const dispatch = useDispatch();

  const handleLogout = () => {
    dispatch(logout());
  };

  const { isAuthenticated, user } = auth;

  const items = [
    {
      label: <Link to="/">Home</Link>,
      key: 'home',
    },
    !isAuthenticated && {
      label: <Link to="/login">Login</Link>,
      key: 'login',
    },
    !isAuthenticated && {
      label: <Link to="/register">Register</Link>,
      key: 'register',
    },
    isAuthenticated && user.role_name === 'Admin' && {
      label: <Link to="/admin/dashboard">Admin Dashboard</Link>,
      key: 'admin',
    },
    isAuthenticated && user.role_name === 'Moderator' && {
      label: <Link to="/moderator/dashboard">Moderator Dashboard</Link>,
      key: 'moderator',
    },
    isAuthenticated && user.role_name === 'Seller' && {
      label: <Link to="/seller/dashboard">Seller Dashboard</Link>,
      key: 'seller',
    },
    isAuthenticated && user.role_name === 'Customer' && {
      label: <Link to="/customer/dashboard">Customer Dashboard</Link>,
      key: 'customer',
    },
    isAuthenticated && {
      label: (
        <span onClick={handleLogout} style={{ cursor: 'pointer' }}>
          Logout
        </span>
      ),
      key: 'logout',
    },
  ].filter(Boolean); // Filter out null or undefined items

  return <Menu mode="horizontal" items={items} />;
}

export default Navbar;
