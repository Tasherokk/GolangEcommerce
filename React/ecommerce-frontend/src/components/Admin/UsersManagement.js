import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, Button, message, Select, Spin } from 'antd';

const { Option } = Select;

function UsersManagement() {
    const [users, setUsers] = useState([]);
    const [roles, setRoles] = useState([]); // Список доступных ролей
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadData();
    }, []);

    // Функция для загрузки данных (пользователей и ролей)
    const loadData = async () => {
        try {
            setLoading(true);
            const [usersResponse, rolesResponse] = await Promise.all([
                axios.get('/api/admin/users'),
                axios.get('/api/admin/roles'),
            ]);
            setUsers(usersResponse.data);
            setRoles(rolesResponse.data);
        } catch (error) {
            console.error('Error loading data:', error);
            message.error('Failed to load data');
        } finally {
            setLoading(false);
        }
    };

    // Обработка удаления пользователя
    const handleDeleteUser = async (userId) => {
        try {
            await axios.delete(`/api/admin/users/${userId}`);
            message.success('User deleted successfully');
            loadData(); // Обновляем данные после удаления
        } catch (error) {
            console.error('Error deleting user:', error);
            message.error('Failed to delete user');
        }
    };

    // Обработка изменения роли
    const handleRoleChange = async (userId, newRole) => {
        try {
            await axios.put(`/api/admin/users/${userId}/role`, { role_name: newRole });
            message.success('User role updated successfully');
            loadData(); // Обновляем данные после изменения роли
        } catch (error) {
            console.error('Error updating user role:', error);
            message.error('Failed to update user role');
        }
    };

    const columns = [
        {
            title: 'User ID',
            dataIndex: 'user_id',
            key: 'user_id',
            sorter: (a, b) => a.user_id - b.user_id, // Сортировка по user_id
        },
        {
            title: 'Username',
            dataIndex: 'username',
            key: 'username',
            sorter: (a, b) => a.username.localeCompare(b.username), // Сортировка по имени
        },
        {
            title: 'Email',
            dataIndex: 'email',
            key: 'email',
            sorter: (a, b) => a.email.localeCompare(b.email), // Сортировка по email
        },
        {
            title: 'Role',
            dataIndex: 'role_name',
            key: 'role_name',
            sorter: (a, b) => a.role_name.localeCompare(b.role_name), // Сортировка по роли
            render: (roleName, record) => (
                <Select
                    defaultValue={roleName}
                    onChange={(newRole) => handleRoleChange(record.user_id, newRole)}
                    style={{ width: 150 }}
                >
                    {roles.map((role) => (
                        <Option key={role.role_id} value={role.role_name}>
                            {role.role_name}
                        </Option>
                    ))}
                </Select>
            ),
        },
        {
            title: 'Action',
            key: 'action',
            render: (_, record) => (
                <Button
                    danger
                    onClick={() => handleDeleteUser(record.user_id)}
                >
                    Delete
                </Button>
            ),
        },
    ];

    return (
        <div>
            <h2>Manage Users</h2>
            {loading ? (
                <Spin size="large" style={{ display: 'block', margin: '50px auto' }} />
            ) : (
                <Table
                    dataSource={users}
                    columns={columns}
                    rowKey="user_id"
                    pagination={{ pageSize: 10 }}
                    bordered
                />
            )}
        </div>
    );
}

export default UsersManagement;
