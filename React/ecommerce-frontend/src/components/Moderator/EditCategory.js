import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Form, Input, Button, message } from 'antd';

function EditCategory() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [form] = Form.useForm();

    useEffect(() => {
        axios.get(`/api/categories/${id}`)
            .then((res) => {
                form.setFieldsValue(res.data);
                setLoading(false);
            })
            .catch((err) => {
                console.error('Error fetching category:', err);
                setLoading(false);
            });
    }, [id, form]);

    const handleUpdate = (values) => {
        axios.put(`/api/categories/${id}`, values)
            .then(() => {
                message.success('Category updated successfully');
                navigate('/moderator/dashboard/categories'); // Возвращаемся к списку категорий
            })
            .catch((err) => {
                console.error('Error updating category:', err);
                message.error('Failed to update category');
            });
    };

    return (
        <div style={{ padding: '20px' }}>
            <h2>Edit Category</h2>
            <Form
                form={form}
                onFinish={handleUpdate}
                layout="vertical"
            >
                <Form.Item name="name" label="Category Name" rules={[{ required: true }]}>
                    <Input />
                </Form.Item>
                <Form.Item name="description" label="Description">
                    <Input.TextArea rows={4} />
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        Save Changes
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default EditCategory;
