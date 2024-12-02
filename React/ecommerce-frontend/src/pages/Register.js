import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { register } from '../store/actions/authActions';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Select, Alert } from 'antd';

const { Option } = Select;

function Register() {
    const dispatch = useDispatch();
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const onFinish = values => {
        dispatch(register(values))
            .then(() => {
                navigate('/'); // Перенаправляем на главную страницу
            })
            .catch(() => {
                setError('Registration failed');
            });
    };

    return (
        <div style={{ maxWidth: 400, margin: '0 auto', padding: '50px 0' }}>
            <h2>Register</h2>
            {error && <Alert message={error} type="error" showIcon />}
            <Form onFinish={onFinish}>
                <Form.Item
                    name="username"
                    rules={[{ required: true, message: 'Please input your username!' }]}
                >
                    <Input placeholder="Username" />
                </Form.Item>
                <Form.Item
                    name="email"
                    rules={[{ required: true, message: 'Please input your email!', type: 'email' }]}
                >
                    <Input placeholder="Email" />
                </Form.Item>
                <Form.Item
                    name="password"
                    rules={[{ required: true, message: 'Please input your password!' }]}
                >
                    <Input.Password placeholder="Password" />
                </Form.Item>
                <Form.Item
                    name="role_name"
                    rules={[{ required: true, message: 'Please select your role!' }]}
                >
                    <Select placeholder="Select a role">
                        <Option value="Customer">Customer</Option>
                        <Option value="Seller">Seller</Option>
                        {/* <Option value="Guest">Guest</Option> */}
                    </Select>
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" block>
                        Register
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default Register;
