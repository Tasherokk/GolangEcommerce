import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table } from 'antd';

function RolesManagement() {
    const [roles, setRoles] = useState([]);

    useEffect(() => {
        fetchRoles();
    }, []);

    const fetchRoles = () => {
        axios.get('/api/admin/roles')
            .then(res => setRoles(res.data))
            .catch(err => console.error(err));
    };

    const columns = [
        { title: 'Role ID', dataIndex: 'role_id', key: 'role_id' },
        { title: 'Role Name', dataIndex: 'role_name', key: 'role_name' },
    ];

    return (
        <div>
            <h2>Manage Roles</h2>
            <Table dataSource={roles} columns={columns} rowKey="role_id" />
        </div>
    );
}

export default RolesManagement;
