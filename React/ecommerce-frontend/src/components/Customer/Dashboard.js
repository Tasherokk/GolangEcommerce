import React from 'react';
import { Link, Route, Routes } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import OrderHistory from './OrderHistory';
import Favorites from './Favorites';
import Profile from '../common/Profile';
import Cart from './Cart';
import OrderDetails from './OrderDetails';

const { Content, Sider } = Layout;

function CustomerDashboard() {

    const menuItems = [
        {
            key: 'cart',
            label: <Link to={`http://localhost:3333/customer/dashboard/cart`}>Cart</Link>,
        },
        {
            key: 'orders',
            label: <Link to={`http://localhost:3333/customer/dashboard/orders`}>Order History</Link>,
        },
        {
            key: 'favorites',
            label: <Link to={`http://localhost:3333/customer/dashboard/favorites`}>Favorites</Link>,
        },
        {
            key: 'profile',
            label: <Link to={`http://localhost:3333/customer/dashboard/profile`}>Profile</Link>,
        },
    ];

    return (
        <Layout>
            <Sider width={200}>
                <Menu mode="inline" defaultSelectedKeys={[]} items={menuItems} />
            </Sider>
            <Layout>
                <Content style={{ padding: '0 24px', minHeight: 280 }}>
                    <Routes>
                        <Route path="orders" element={<OrderHistory />} />
                        <Route path="orders/:id" element={<OrderDetails />} />
                        <Route path="favorites" element={<Favorites />} />
                        <Route path="profile" element={<Profile />} />
                        <Route path="cart" element={<Cart />} />
                        {/* Главная страница панели */}
                        <Route
                            path="/"
                            element={<h1>Welcome to Customer Dashboard</h1>}
                        />
                    </Routes>
                </Content>
            </Layout>
        </Layout>
    );
}

export default CustomerDashboard;
