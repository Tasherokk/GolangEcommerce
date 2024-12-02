import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, Tag, Typography } from 'antd';

const { Title } = Typography;

function MyOrders() {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchOrders();
    }, []);

    const fetchOrders = () => {
        setLoading(true);
        axios
            .get('/api/seller/orders') // Эндпоинт для получения заказов
            .then((res) => {
                setOrders(res.data);
                setLoading(false);
            })
            .catch((err) => {
                console.error('Error fetching orders:', err);
                setLoading(false);
            });
    };

    const columns = [
        {
            title: 'Order ID',
            dataIndex: 'order_id',
            key: 'order_id',
        },
        {
            title: 'Product Name',
            dataIndex: 'product_name',
            key: 'product_name',
        },
        {
            title: 'Quantity',
            dataIndex: 'quantity',
            key: 'quantity',
        },
        {
            title: 'Price',
            dataIndex: 'price',
            key: 'price',
            render: (price) => `$${price.toFixed(2)}`,
        },
        {
            title: 'Total Amount',
            dataIndex: 'total_amount',
            key: 'total_amount',
            render: (totalAmount) => `$${totalAmount.toFixed(2)}`,
        },
        {
            title: 'Order Date',
            dataIndex: 'order_date',
            key: 'order_date',
            render: (date) => new Date(date).toLocaleDateString(),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => (
                <Tag color={status === 'Paid' ? 'green' : 'red'}>
                    {status.toUpperCase()}
                </Tag>
            ),
        },
    ];

    return (
        <div>
            <Title level={2}>My Orders</Title>
            <Table
                dataSource={orders}
                columns={columns}
                rowKey="order_item_id"
                loading={loading}
                pagination={{ pageSize: 10 }}
            />
        </div>
    );
}

export default MyOrders;
