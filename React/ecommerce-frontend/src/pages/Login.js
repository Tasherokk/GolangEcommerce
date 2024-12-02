import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { login } from '../store/actions/authActions';
import { Form, Input, Button, Alert } from 'antd';
import { useNavigate } from 'react-router-dom';

function Login() {
    const dispatch = useDispatch();
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const onFinish = values => {
        dispatch(login(values))
            .then(() => {
                navigate('/'); // Перенаправление на homepage
            })
            .catch(() => {
                setError('Invalid credentials');
            });
    };

    return (
        <div style={{ maxWidth: 400, margin: '0 auto', padding: '50px 0' }}>
            <h2>Login</h2>
            {error && <Alert message={error} type="error" showIcon />}
            <Form onFinish={onFinish}>
                <Form.Item
                    name="username"
                    rules={[{ required: true, message: 'Please input your username!' }]}
                >
                    <Input placeholder="Username" />
                </Form.Item>
                <Form.Item
                    name="password"
                    rules={[{ required: true, message: 'Please input your password!' }]}
                >
                    <Input.Password placeholder="Password" />
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" block>
                        Login
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default Login;
