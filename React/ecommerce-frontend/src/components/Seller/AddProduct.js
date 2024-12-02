import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Form, Input, InputNumber, Button, Select, message } from 'antd';

const { Option } = Select;

function AddProduct() {
    const [form] = Form.useForm();
    const [categories, setCategories] = useState([]);

    // Загрузка категорий
    useEffect(() => {
        axios.get('/api/categories')
            .then(res => setCategories(res.data))
            .catch(err => console.error('Error fetching categories:', err));
    }, []);

    const onFinish = (values) => {
        axios.post('/api/products', values)
            .then((res) => {
                message.success(`Product added successfully with ID: ${res.data.product_id}`);
                form.resetFields();
            })
            .catch((err) => {
                console.error('Error adding product:', err);
                message.error('Failed to add product');
            });
    };

    return (
        <Form form={form} onFinish={onFinish} layout="vertical">
            <Form.Item
                name="name"
                label="Product Name"
                rules={[{ required: true, message: 'Please enter the product name' }]}
            >
                <Input />
            </Form.Item>
            <Form.Item
                name="description"
                label="Description"
            >
                <Input.TextArea />
            </Form.Item>
            <Form.Item
                name="price"
                label="Price"
                rules={[{ required: true, message: 'Please enter the price' }]}
            >
                <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item
                name="stock"
                label="Stock"
                rules={[{ required: true, message: 'Please enter the stock quantity' }]}
            >
                <InputNumber min={0} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item
                name="category_id"
                label="Category"
                rules={[{ required: true, message: 'Please select a category' }]}
            >
                <Select placeholder="Select a category">
                    {categories.map(category => (
                        <Option key={category.category_id} value={category.category_id}>
                            {category.name}
                        </Option>
                    ))}
                </Select>
            </Form.Item>
            <Form.Item>
                <Button type="primary" htmlType="submit">Add Product</Button>
            </Form.Item>
        </Form>
    );
}

export default AddProduct;
