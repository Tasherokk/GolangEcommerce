import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, Button, message } from 'antd';
import { Link } from 'react-router-dom';

function MyProducts() {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchMyProducts();
    }, []);

    const fetchMyProducts = () => {
        setLoading(true);
        axios.get('/api/seller/products')
            .then(res => {
                setProducts(res.data);
                setLoading(false);
            })
            .catch(err => {
                console.error(err);
                setLoading(false);
            });
    };

    const handleDelete = (productId) => {
        axios.delete(`/api/products/${productId}`)
            .then(() => {
                message.success('Product deleted');
                fetchMyProducts();
            })
            .catch(err => {
                console.error(err);
                message.error('Failed to delete product');
            });
    };

    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            sorter: (a, b) => a.name.localeCompare(b.name), // Сортировка по имени
            render: (text, record) => <Link to={`/seller/dashboard/my-products/${record.product_id}/edit`}>{text}</Link>, // Ссылка на ProductEdit
        },
        {
            title: 'Price',
            dataIndex: 'price',
            key: 'price',
            sorter: (a, b) => a.price - b.price, // Сортировка по цене
            render: (price) => `$${price.toFixed(2)}`, // Отображение с двумя знаками после запятой
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
                <>
                    <Button
                        type="link"
                        onClick={() => handleDelete(record.product_id)}
                        danger
                    >
                        Delete
                    </Button>
                </>
            ),
        },
    ];

    return (
        <div>
            <h2>My Products</h2>
            <Table
                dataSource={products}
                columns={columns}
                rowKey="product_id"
                loading={loading}
                pagination={{ pageSize: 10 }} // Пагинация по 10 продуктов на странице
            />
        </div>
    );
}

export default MyProducts;
