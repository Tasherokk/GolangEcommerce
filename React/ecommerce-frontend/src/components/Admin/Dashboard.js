import React from 'react';
import { Link, Route, Routes } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import UsersManagement from './UsersManagement';
import ProductsManagement from './ProductsManagement';
import CategoriesManagement from './CategoriesManagement';
import RolesManagement from './RolesManagement';

const { Content, Sider } = Layout;

function AdminDashboard() {

    const menuItems = [
        {
            key: 'users',
            label: <Link to={`http://localhost:3333/admin/dashboard/users`}>Manage Users</Link>,
        },
        {
            key: 'roles',
            label: <Link to={`http://localhost:3333/admin/dashboard/roles`}>Manage Roles</Link>,
        },
        {
            key: 'products',
            label: <Link to={`http://localhost:3333/admin/dashboard/products`}>Manage Products</Link>,
        },
        {
            key: 'categories',
            label: <Link to={`http://localhost:3333/admin/dashboard/categories`}>Manage Categories</Link>,
        },
        // Добавьте другие элементы меню
    ];

    return (
        <Layout>
            <Sider width={200}>
                <Menu mode="inline" defaultSelectedKeys={[]} items={menuItems} />
            </Sider>
            <Layout>
                <Content style={{ padding: '0 24px', minHeight: 280 }}>
                    <Routes>
                        <Route path="users" element={<UsersManagement />} />
                        <Route path="products" element={<ProductsManagement />} />
                        <Route path="categories" element={<CategoriesManagement />} />
                        <Route path="roles" element={<RolesManagement />} />
                        {/* Другие маршруты */}
                        <Route path="/" element={<h1>Welcome to Admin Dashboard</h1>} />
                    </Routes>
                </Content>
            </Layout>
        </Layout>
    );
}

export default AdminDashboard;
