import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { Form, Input, InputNumber, Button, message, Typography, Select, Spin } from 'antd';

const { Title } = Typography;
const { Option } = Select;

function ProductsEdit() {
    const { id } = useParams();
    // eslint-disable-next-line no-unused-vars
    const [product, setProduct] = useState(null);
    const [categories, setCategories] = useState([]); // Для загрузки списка категорий
    const [loading, setLoading] = useState(true);
    const [form] = Form.useForm();

    useEffect(() => {
        // Загружаем данные продукта
        axios.get(`/api/products/${id}`)
            .then(res => {
                setProduct(res.data.product);
                form.setFieldsValue({
                    ...res.data.product,
                    category_name: res.data.product.category_name, // Устанавливаем category_name в форму
                });
                setLoading(false);
            })
            .catch(err => {
                console.error('Error fetching product:', err);
                setLoading(false);
            });

        // Загружаем категории
        axios.get('/api/categories')
            .then(res => {
                setCategories(res.data);
            })
            .catch(err => {
                console.error('Error fetching categories:', err);
            });
    }, [id, form]);

    const handleUpdate = (values) => {
        axios.put(`/api/products/${id}`, values)
            .then(() => {
                message.success('Product updated successfully');
            })
            .catch(err => {
                console.error('Error updating product:', err);
                message.error('Failed to update product');
            });
    };

    return (
        <div style={{ padding: '20px' }}>
            <Spin spinning={loading} tip="Loading..." size="large">
                {!loading && (
                    <div>
                        <Title level={2}>Edit Product</Title>
                        <Form
                            form={form}
                            onFinish={handleUpdate}
                            layout="vertical"
                        >
                            <Form.Item name="name" label="Product Name" rules={[{ required: true }]}>
                                <Input />
                            </Form.Item>
                            <Form.Item name="description" label="Description">
                                <Input.TextArea rows={4} />
                            </Form.Item>
                            <Form.Item name="price" label="Price" rules={[{ required: true }]}>
                                <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
                            </Form.Item>
                            <Form.Item name="stock" label="Stock" rules={[{ required: true }]}>
                                <InputNumber min={0} style={{ width: '100%' }} />
                            </Form.Item>
                            <Form.Item name="category_name" label="Category" rules={[{ required: true }]}>
                                <Select placeholder="Select a category">
                                    {categories.map((category) => (
                                        <Option key={category.category_id} value={category.name}>
                                            {category.name}
                                        </Option>
                                    ))}
                                </Select>
                            </Form.Item>
                            <Form.Item name="image_url" label="Image URL">
                                <Input />
                            </Form.Item>
                            <Form.Item>
                                <Button type="primary" htmlType="submit">
                                    Save Changes
                                </Button>
                            </Form.Item>
                        </Form>
                    </div>
                )}
            </Spin>
        </div>
    );
}

export default ProductsEdit;
