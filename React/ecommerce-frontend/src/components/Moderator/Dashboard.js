import React from 'react';
import { Link, Routes, Route } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import ProductsModeration from './ProductsModeration';
import CategoriesModeration from './CategoriesModeration';
import EditCategory from './EditCategory';

const { Content, Sider } = Layout;

function ModeratorDashboard() {
    const menuItems = [
        {
            key: 'products',
            label: <Link to="/moderator/dashboard/products">Moderate Products</Link>,
        },
        {
            key: 'categories',
            label: <Link to="/moderator/dashboard/categories">Moderate Categories</Link>,
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
                        <Route path="products" element={<ProductsModeration />} />
                        <Route path="categories" element={<CategoriesModeration />} />
                        <Route path="categories/:id/edit" element={<EditCategory />} />
                        <Route path="/" element={<h1>Welcome to Moderator Dashboard</h1>} />
                    </Routes>
                </Content>
            </Layout>
        </Layout>
    );
}

export default ModeratorDashboard;
