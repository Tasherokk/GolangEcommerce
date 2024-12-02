import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, Button, message } from 'antd';

function ProductsManagement() {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchProducts = async () => {
        try {
            setLoading(true);
            const res = await axios.get('/api/products');
            setProducts(res.data);
        } catch (err) {
            console.error('Error fetching products:', err);
            message.error('Failed to load products');
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteProduct = async (productId) => {
        try {
            await axios.delete(`/api/products/${productId}`);
            message.success('Product deleted successfully');
            fetchProducts(); // Обновляем данные после удаления
        } catch (err) {
            console.error('Error deleting product:', err);
            message.error('Failed to delete product');
        }
    };

    const columns = [
        {
            title: 'Product ID',
            dataIndex: 'product_id',
            key: 'product_id',
            sorter: (a, b) => a.product_id - b.product_id, // Сортировка по ID продукта
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            sorter: (a, b) => a.name.localeCompare(b.name), // Сортировка по имени продукта
        },
        {
            title: 'Description',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: 'Price',
            dataIndex: 'price',
            key: 'price',
            sorter: (a, b) => a.price - b.price, // Сортировка по цене
        },
        {
            title: 'Stock',
            dataIndex: 'stock',
            key: 'stock',
            sorter: (a, b) => a.stock - b.stock, // Сортировка по количеству на складе
        },
        {
            title: 'Action',
            key: 'action',
            render: (_, record) => (
                <Button danger onClick={() => handleDeleteProduct(record.product_id)}>
                    Delete
                </Button>
            ),
        },
    ];

    return (
        <div>
            <h2>Manage Products</h2>
            <Table
                dataSource={products}
                columns={columns}
                rowKey="product_id"
                loading={loading}
                pagination={{ pageSize: 10 }} // Добавлена пагинация
                bordered
            />
        </div>
    );
}

export default ProductsManagement;
