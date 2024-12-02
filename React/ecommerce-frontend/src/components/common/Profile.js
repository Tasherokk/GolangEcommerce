import React, { useEffect, useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import axios from 'axios';

function Profile() {
  const [user, setUser] = useState({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Загрузка информации о пользователе
    axios
      .get('/api/users/me') // Эндпоинт для получения текущего пользователя
      .then((res) => {
        setUser(res.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Error loading user data:', err);
        message.error('Failed to load user data');
        setLoading(false);
      });
  }, []);

  const handleFinish = (values) => {
    // Отправка обновлений профиля
    axios
      .put('/api/users/me', values)
      .then(() => {
        message.success('Profile updated successfully');
        setUser({ ...user, ...values }); // Обновление состояния после успешного обновления
      })
      .catch((err) => {
        console.error('Error updating profile:', err);
        message.error('Failed to update profile');
      });
  };

  if (loading) {
    return <p>Loading profile...</p>;
  }

  return (
    <div>
      <h2>Your Profile</h2>
      <Form
        layout="vertical"
        initialValues={user}
        onFinish={handleFinish}
        style={{ maxWidth: 600, margin: '0 auto' }}
      >
        <Form.Item
          label="Username"
          name="username"
          rules={[{ required: true, message: 'Please enter your username!' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Email"
          name="email"
          rules={[
            { required: true, message: 'Please enter your email!' },
            { type: 'email', message: 'Please enter a valid email!' },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Role"
          name="role_name"
        >
          <Input disabled />

        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Save Changes
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}

export default Profile;
