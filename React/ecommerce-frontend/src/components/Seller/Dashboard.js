import React from 'react';
import { Link, Routes, Route } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import AddProduct from './AddProduct';
import MyProducts from './MyProducts';
import MyOrders from './MyOrders';
import Profile from '../common/Profile';
import ProductsEdit from './ProductsEdit';

const { Content, Sider } = Layout;

function SellerDashboard() {
    const menuItems = [
        {
            key: 'add-product',
            label: <Link to="/seller/dashboard/add-product">Add Product</Link>,
        },
        {
            key: 'my-products',
            label: <Link to="/seller/dashboard/my-products">My Products</Link>,
        },
        {
            key: 'my-orders',
            label: <Link to="/seller/dashboard/my-orders">My Orders</Link>,
        },
        {
            key: 'profile',
            label: <Link to="/seller/dashboard/profile">Profile</Link>,
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
                        <Route path="add-product" element={<AddProduct />} />
                        <Route path="my-products" element={<MyProducts />} />
                        <Route path="my-products/:id/edit" element={<ProductsEdit />} />
                        <Route path="my-orders" element={<MyOrders />} />
                        <Route path="profile" element={<Profile />} />
                        <Route path="/" element={<h1>Welcome to Seller Dashboard</h1>} />
                    </Routes>
                </Content>
            </Layout>
        </Layout>
    );
}

export default SellerDashboard;
